package chat

import (
	"context"
	"github.com/AntNoHuabei/Remo/pkg/api/response"
	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
)

func NewContinuousAgent(ctx context.Context) (*ContinuousAgent, error) {

	cm, err := deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
		APIKey:  "sk-0ed0d69fac2c468ba57ac31fc2faf1e3",
		Model:   "deepseek-reasoner",
		BaseURL: "https://api.deepseek.com/beta",
	})
	if err != nil {
		return nil, err
	}

	agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:          "ContinuousAgent",
		Description:   "I can keep talking with you and remember everything you've said",
		Instruction:   "",
		Model:         cm,
		ToolsConfig:   adk.ToolsConfig{},
		GenModelInput: nil,
		Exit:          nil,
		OutputKey:     "",
		MaxIterations: 0,
		Middlewares:   nil,
	})

	if err != nil {
		return nil, err
	}

	return &ContinuousAgent{
		agent: agent,
	}, nil
}

type ContinuousAgent struct {
	agent    *adk.ChatModelAgent
	session  string
	runner   *adk.Runner
	messages []adk.Message
}

// Recover 从断点恢复
func (agent *ContinuousAgent) Recover(ctx context.Context, session string) error {

	s, err := NewStore(session)
	if err != nil {
		return err
	}
	agent.runner = adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           agent.agent,
		EnableStreaming: true,
		CheckPointStore: s,
	})

	messages, err := Messages(session)

	if err != nil {
		return err
	}

	agent.session = session
	agent.messages = make([]adk.Message, 0)
	for _, message := range messages {
		agent.messages = append(agent.messages, &schema.Message{
			Content: message.Content,
			Role:    schema.RoleType(message.Role),
		})
	}

	return nil
}

func (agent *ContinuousAgent) Abort() {

}

func (agent *ContinuousAgent) Chat(ctx context.Context, message *Message) (<-chan response.ChatResponse, error) {

	if message.RequestId == "" {
		message.RequestId = uuid.New().String()
	}

	//TODO 判断是否有工具调用确认

	agent.messages = append(agent.messages, &schema.Message{
		Content: message.Content,
		Role:    schema.User,
	})

	MessageAppend(agent.session, message)

	ch := make(chan response.ChatResponse)

	it := agent.runner.Run(ctx, agent.messages, adk.WithCheckPointID("session-"+message.RequestId))

	go func() {

		var outputMessage = &schema.Message{
			Role:    schema.Assistant,
			Content: "",
		}
		for {
			event, ok := it.Next()
			if !ok {
				break
			}

			if event.Err != nil {
				ch <- response.ChatResponse{
					Error:     event.Err,
					RequestID: message.RequestId,
				}
			} else {

				if event.Output != nil && event.Output.MessageOutput != nil {

					if event.Output.MessageOutput.MessageStream != nil {

						for {
							m, err := event.Output.MessageOutput.MessageStream.Recv()
							if err != nil {
								break
							}

							if m.ReasoningContent != "" {
								outputMessage.ReasoningContent = m.ReasoningContent
								ch <- response.ChatResponse{
									ReasonContent: m.ReasoningContent,
									RequestID:     message.RequestId,
								}
							} else {
								outputMessage.Content += m.Content
								ch <- response.ChatResponse{
									Content:   m.Content,
									RequestID: message.RequestId,
								}
							}
						}
					}
				}
			}
		}

		agent.messages = append(agent.messages, outputMessage)

		//TODO 存储会话

		MessageAppend(agent.session, &Message{
			Content:   outputMessage.Content,
			Role:      "assistant",
			RequestId: message.RequestId,
		})
		close(ch)
	}()

	return ch, nil

}
