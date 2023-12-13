package main

import (
    "encoding/json"
    "fmt"
    "os"
    "strings"

    process "go-chat/gen"
)

type MessageArchive map[string]string

func parseAddress(addressString string) (node, process_name, package_name, publisher string) {
    parts := strings.Split(addressString, "@")
    node = parts[0]
    parts = strings.Split(parts[1], ":")
    return node, parts[0], parts[1], parts[2]
}

func handleMessage(ourNode string, messageArchive MessageArchive) MessageArchive {
    result := process.UqbarProcess0_4_0_StandardReceive().Unwrap()
    source := result.F0
    message := result.F1
    // source, message, err := process.UqbarProcess0_4_0_StandardReceive()
    // if err != nil {
    //     process.UqbarProcess0_4_0_StandardPrintToTerminal(0, fmt.Sprintf("Error: %v", err))
    //     os.Exit(1)
    // }

    switch message.Kind() {
    case process.UqbarProcess0_4_0_StandardMessageKindResponse:
        process.UqbarProcess0_4_0_StandardPrintToTerminal(
            0,
            fmt.Sprintf("chat: unexpected Response: %v\n", message),
        )
        os.Exit(1)
    case process.UqbarProcess0_4_0_StandardMessageKindRequest:
        var ipc map[string]interface{}
        err := json.Unmarshal(message.GetRequest().Ipc, &ipc)
        if err != nil {
            fmt.Println("Error:", err)
            os.Exit(1)
        }

        if send, ok := ipc["Send"].(map[string]interface{}); ok {
            target, messageText := send["target"].(string), send["message"].(string)
            if target == ourNode {
                process.UqbarProcess0_4_0_StandardPrintToTerminal(
                    0,
                    fmt.Sprintf("chat|%s: %s\n", source.Node, messageText),
                )
                messageArchive[source.Node] = messageText
            } else {
                process.UqbarProcess0_4_0_StandardSendAndAwaitResponse(
                    process.UqbarProcess0_4_0_StandardAddress{
                        Node: target,
                        Process: process.UqbarProcess0_4_0_StandardProcessId{
                            ProcessName: "chat",
                            PackageName: "chat",
                            PublisherNode: "uqbar",
                        },
                    },
                    process.UqbarProcess0_4_0_StandardRequest{
                        Inherit: false,
                        ExpectsResponse: process.Some(uint64(5)),
                        Ipc: message.GetRequest().Ipc,
                        Metadata: process.None[process.UqbarProcess0_4_0_StandardJson](),
                    },
                    process.None[process.UqbarProcess0_4_0_StandardPayload](),
                )
            }
            response := process.UqbarProcess0_4_0_StandardResponse{
                Inherit: false,
                Ipc: []byte(`{"Ack": null}`),
                Metadata: process.None[process.UqbarProcess0_4_0_StandardJson](),
            }
            process.UqbarProcess0_4_0_StandardSendResponse(
                response,
                process.None[process.UqbarProcess0_4_0_StandardPayload](),
            )
        } else if _, ok := ipc["History"].(map[string]interface{}); ok {
            response := process.UqbarProcess0_4_0_StandardResponse{
                Inherit: false,
                Ipc: []byte(fmt.Sprintf(`{"History": {"messages": %v}}`, messageArchive)),
                Metadata: process.None[process.UqbarProcess0_4_0_StandardJson](),
            }
            process.UqbarProcess0_4_0_StandardSendResponse(
                response,
                process.None[process.UqbarProcess0_4_0_StandardPayload](),
            )
        } else {
            os.Exit(1)
        }
    }

    return messageArchive
}

func init() {
    a := ProcessImpl{}
    process.SetProcess(a)
}

type ProcessImpl struct {}

func (p ProcessImpl) Init(our string) {
    process.UqbarProcess0_4_0_StandardPrintToTerminal(
        0,
        "chat: begin (golang)",
    )

    ourNode, _, _, _ := parseAddress(our)

    messageArchive := make(MessageArchive)

    for {
        messageArchive = handleMessage(ourNode, messageArchive)
    }
}

//go:generate wit-bindgen tiny-go wit/ -w process --out-dir=gen
func main() {}
