package rpc_server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"participle/internal"
	"participle/logger"
)

type Reply struct {
}

// === Article Classification ===

type ArticleReqMsg struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type ArticleReplyMsg struct {
	Code    int          `json:"code"`
	Content articleReply `json:"content"`
}

type articleReply struct {
	LogID int             `json:"log_id"`
	Items articleItemList `json:"item"`
}

type articleItems struct {
	Score float64 `json:"score"`
	Tag   string  `json:"tag"`
}

type articleItemList struct {
	Lv2TagList []articleItems `json:"lv2_tag_list"`
	Lv1TagList []articleItems `json:"lv1_tag_list"`
}

// === Emotion analysis ===

type EmotionReqMsg struct {
	Text string `json:"text"`
}

type EmotionReplyMsg struct {
	Code    int          `json:"code"`
	Content emotionReply `json:"content"`
}

type emotionReply struct {
	Text  string         `json:"text"`
	Items []emotionItems `json:"items"`
}
type emotionItems struct {
	Sentiment    int     `json:"sentiment"`
	Confidence   float64 `json:"confidence"`
	PositiveProb float64 `json:"positive_prob"`
	NegativeProb float64 `json:"negative_prob"`
}

// === Lexer ===

type LexerReqMsg struct {
	Text string `json:"text"`
}

type LexerReplyMsg struct {
	Code    int        `json:"code"`
	Content lexerReply `json:"content"`
}

type lexerReply struct {
	Text  string       `json:"text"`
	Items []lexerItems `json:"items"`
}
type lexerItems struct {
	ByteLength int           `json:"byte_length"`
	ByteOffset int           `json:"byte_offset"`
	Formal     string        `json:"formal"`
	Item       string        `json:"item"`
	Ne         string        `json:"ne"`
	Pos        string        `json:"pos"`
	URI        string        `json:"uri"`
	LocDetails []localDetail `json:"loc_details"`
	BasicWords []string      `json:"basic_words"`
}

type localDetail struct {
	Type       string `json:"type"`
	ByteOffset int    `json:"byte_offset"`
	ByteLength int    `json:"byte_length"`
}

type connection struct {
	io.Writer
	io.ReadCloser
}

func (r *Reply) ArticleClassification(req ArticleReqMsg, reply *ArticleReplyMsg) error {
	title, content := req.Title, req.Content
	titleByte, contentByte := []byte(title), []byte(content)
	if len(titleByte) > 80 || len(contentByte) > 65535 {
		logger.Log.Error("Article title or content is too long!")
		*reply = ArticleReplyMsg{Code: 0}
		return nil
	}

	token, err := internal.GetToken()
	if err != nil {
		logger.Log.Error(fmt.Sprintf("[Article] get token err: %s", err.Error()))
		*reply = ArticleReplyMsg{Code: 0}
		return nil
	}

	reqContentByte, err := json.Marshal(&req)
	if err != nil {
		logger.Log.Error(
			fmt.Sprintf("Failed to json marshal article reqContent, err: %s", err.Error()))
		*reply = ArticleReplyMsg{Code: 0}
		return nil
	}

	articleUrl := fmt.Sprintf(
		"https://aip.baidubce.com/rpc/2.0/nlp/v1/topic?charset=UTF-8&access_token=%s",
		token)

	articleReq, err := http.NewRequest("POST", articleUrl, bytes.NewReader(reqContentByte))
	if err != nil {
		logger.Log.Error("Failed to create new article request")
		*reply = ArticleReplyMsg{Code: 0}
		return nil
	}

	articleReq.Header.Set("Content-Type", "application/json")

	proxyClient := &http.Client{}

	articleRes, err := proxyClient.Do(articleReq)
	if err != nil {
		logger.Log.Error("Failed to get baidu article response!")
		*reply = ArticleReplyMsg{Code: 0}
		return nil
	}

	defer articleRes.Body.Close()

	articleResContent, err := ioutil.ReadAll(articleRes.Body)
	if err != nil {
		logger.Log.Error("Failed to read baidu response body!")
		*reply = ArticleReplyMsg{Code: 0}
		return nil
	}

	articleReplyContent := articleReply{}

	err = json.Unmarshal(articleResContent, &articleReplyContent)
	if err != nil {
		logger.Log.Error("Failed to unmarshal baidu reply content!")
		*reply = ArticleReplyMsg{Code: 0}
		return nil
	}

	*reply = ArticleReplyMsg{Code: 1, Content: articleReplyContent}

	return nil
}

func (r *Reply) EmotionAnalysis(req EmotionReqMsg, reply *EmotionReplyMsg) error {
	textByte := []byte(req.Text)
	if len(textByte) > 2048 {
		logger.Log.Error("Emotion text too long!")
		*reply = EmotionReplyMsg{Code: 0}
		return nil
	}

	reqContentByte, err := json.Marshal(&req)
	if err != nil {
		logger.Log.Error(
			fmt.Sprintf("Failed to json marshal emotion reqContent, err: %s", err.Error()))
		*reply = EmotionReplyMsg{Code: 0}
		return nil
	}

	token, err := internal.GetToken()
	if err != nil {
		logger.Log.Error(fmt.Sprintf("[Emotion] get token err: %s", err.Error()))
		*reply = EmotionReplyMsg{Code: 0}
		return nil
	}

	emotionUrl := fmt.Sprintf(
		"https://aip.baidubce.com/rpc/2.0/nlp/v1/sentiment_classify?charset=UTF-8&access_token=%s",
		token)

	emotionReq, err := http.NewRequest("POST", emotionUrl, bytes.NewReader(reqContentByte))
	if err != nil {
		logger.Log.Error("Failed to create new emotion request")
		*reply = EmotionReplyMsg{Code: 0}
		return nil
	}

	emotionReq.Header.Set("Content-Type", "application/json")

	proxyClient := &http.Client{}

	emotionRes, err := proxyClient.Do(emotionReq)
	if err != nil {
		logger.Log.Error("Failed to get baidu emotion response!")
		*reply = EmotionReplyMsg{Code: 0}
		return nil
	}

	defer emotionRes.Body.Close()

	emotionResContent, err := ioutil.ReadAll(emotionRes.Body)
	if err != nil {
		logger.Log.Error("Failed to read baidu response body!")
		*reply = EmotionReplyMsg{Code: 0}
		return nil
	}

	emotionReplyContent := emotionReply{}

	err = json.Unmarshal(emotionResContent, &emotionReplyContent)
	if err != nil {
		logger.Log.Error("Failed to unmarshal baidu reply content!")
		*reply = EmotionReplyMsg{Code: 0}
		return nil
	}

	*reply = EmotionReplyMsg{Code: 1, Content: emotionReplyContent}

	return nil

}

func (r *Reply) Lexer(req LexerReqMsg, reply *LexerReplyMsg) error {
	textByte := []byte(req.Text)
	if len(textByte) > 20000 {
		logger.Log.Error("Lexer text is too long!")
		*reply = LexerReplyMsg{Code: 0}
		return nil
	}

	token, err := internal.GetToken()
	if err != nil {
		logger.Log.Error(fmt.Sprintf("[Lexer] get token err: %s", err.Error()))
		*reply = LexerReplyMsg{Code: 0}
		return nil
	}

	reqContentByte, err := json.Marshal(&req)
	if err != nil {
		logger.Log.Error(
			fmt.Sprintf("Failed to json marshal lexer reqContent, err: %s", err.Error()))
		*reply = LexerReplyMsg{Code: 0}
		return nil
	}

	lexerUrl := fmt.Sprintf(
		"https://aip.baidubce.com/rpc/2.0/nlp/v1/lexer?charset=UTF-8&access_token=%s",
		token)

	lexerReq, err := http.NewRequest("POST", lexerUrl, bytes.NewReader(reqContentByte))
	if err != nil {
		logger.Log.Error("Failed to create new lexer request")
		*reply = LexerReplyMsg{Code: 0}
		return nil
	}

	lexerReq.Header.Set("Content-Type", "application/json")

	proxyClient := &http.Client{}

	lexerRes, err := proxyClient.Do(lexerReq)
	if err != nil {
		logger.Log.Error("Failed to get baidu lexer response!")
		*reply = LexerReplyMsg{Code: 0}
		return nil
	}

	defer lexerRes.Body.Close()

	lexerResContent, err := ioutil.ReadAll(lexerRes.Body)
	if err != nil {
		logger.Log.Error("Failed to read baidu response body!")
		*reply = LexerReplyMsg{Code: 0}
		return nil
	}

	lexerReplyContent := lexerReply{}

	err = json.Unmarshal(lexerResContent, &lexerReplyContent)
	if err != nil {
		logger.Log.Error("Failed to unmarshal baidu reply content!")
		*reply = LexerReplyMsg{Code: 0}
		return nil
	}

	*reply = LexerReplyMsg{Code: 1, Content: lexerReplyContent}

	return nil
}
