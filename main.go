package main

import (
	"github.com/ReneKroon/ttlcache/v2"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"os"
	"strings"
	"time"
)


var banDuration = time.Minute * 1
var cfg *Config

func main() {

	conf,err := ConfigFromFile(os.Args[1])
	cfg = conf
	if err != nil {
		log.Fatal(err)
	}
	b, err := tb.NewBot(tb.Settings{
		Token: cfg.Server.Token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	defer b.Start()
	newItemCallback := func(key string, value interface{}) {
		log.Printf("User(%s) added to ban\n", key)
		message := value.(*tb.Message)
		b.Reply(message,message.Sender.Username+" T'HO BANNATO PE 1 MINUTO!")
	}

	expirationCallback := func(key string, value interface{}) {
		log.Printf("User(%s) removed from ban\n", key)
		message := value.(*tb.Message)
		b.Reply(message,message.Sender.Username+" SEI STATO SBANNATO!")
	}

	cache := ttlcache.NewCache()
	cache.SetTTL(banDuration)
	cache.SetExpirationCallback(expirationCallback)
	cache.SetNewItemCallback(newItemCallback)
	cache.SkipTTLExtensionOnHit(true)

	if err != nil {
		log.Fatal(err)
		return
	}


	b.Handle(tb.OnText, func(m *tb.Message) {
		log.Println(m.Sender.Username+" : "+m.Text)
		present,_ := cache.Get(m.Sender.Username)
		if present != nil {
			b.Delete(m)
		}else if checkBanTerms(m.Text) {
			cache.Set(m.Sender.Username,(m))
		}
	})

	b.Handle(tb.OnVoice,func(m *tb.Message) {
		deleteMsgIfBanned(m,b,cache)
	})
	b.Handle(tb.OnAnimation,func(m *tb.Message) {
		deleteMsgIfBanned(m,b,cache)
	})
	b.Handle(tb.OnPhoto,func(m *tb.Message) {
		deleteMsgIfBanned(m,b,cache)
	})
	b.Handle(tb.OnAudio,func(m *tb.Message) {
		deleteMsgIfBanned(m,b,cache)
	})
	b.Handle(tb.OnVideo,func(m *tb.Message) {
		deleteMsgIfBanned(m,b,cache)
	})
	b.Handle(tb.OnVideoNote,func(m *tb.Message) {
		deleteMsgIfBanned(m,b,cache)
	})


}

func deleteMsgIfBanned(m *tb.Message,b *tb.Bot, cache *ttlcache.Cache){
	present,_ := cache.Get(m.Sender.Username)
	if present != nil {
		b.Delete(m)
	}
}

func checkBanTerms(message string) bool{
	bannedTerms := cfg.Server.Banterms
	lowermsg := strings.ToLower(message)
	for i := range bannedTerms {
		if strings.Contains(lowermsg,bannedTerms[i]) {
			return true
		}
	}
	return false
}
