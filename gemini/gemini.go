package gemini

import (
	"context"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"log"
	"log/slog"
	"net/http"
	"net/url"
)

type (
	GeminiWorker struct {
	}
)

func NewGeminiWorker() GeminiWorker {
	return GeminiWorker{}
}

func (g GeminiWorker) Do() {
	g.do()
}

func (g GeminiWorker) do() {
	ctx := context.Background()
	proxyURL, err := url.Parse("http://proxy.example.com:8080")
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	// 创建自定义的 Client
	hpClient := &http.Client{
		Transport: transport,
	}
	client, err := genai.NewClient(ctx, option.WithAPIKey("AIzaSyDptTcMLyRHIp-xgehxe2CNJVaI4u8RwvU"), option.WithHTTPClient(hpClient))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-pro")
	resp, err := model.GenerateContent(ctx, genai.Text("Write a story about a magic backpack."))
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("resp", resp)
}
