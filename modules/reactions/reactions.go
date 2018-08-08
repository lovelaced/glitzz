package reactions

import (
	"errors"
	"github.com/lovelaced/glitzz/config"
	"github.com/lovelaced/glitzz/core"
	"github.com/lovelaced/glitzz/util"
	"strings"
)

var cutelist = []string{
	"✿◕ ‿ ◕✿",
	"❀◕ ‿ ◕❀",
	"(✿◠‿◠)",
	"(◕‿◕✿) ",
	"( ｡◕‿◕｡)",
	"(◡‿◡✿)",
	"⊂◉‿◉つ ❤",
	"{ ◕ ◡ ◕}",
	"( ´・‿-) ~ ♥",
	"(っ⌒‿⌒)っ~ ♥",
	"ʕ´•ᴥ•`ʔσ”",
	"(･Θ･) caw",
	"(=^･ω･^)y＝",
	"ヽ(=^･ω･^=)丿",
	"~(=^･ω･^)ヾ(^^ )",
	"| (•□•) | (❍ᴥ❍ʋ)",
	"ϞϞ(๑⚈ ․̫ ⚈๑)∩",
	"ヾ(･ω･*)ﾉ",
	"▽・ω・▽ woof~",
	"(◎｀・ω・´)人(´・ω・｀*)",
	"(*´・ω・)ノ(-ω-｀*)",
	"(❁´ω`❁)",
	"(＊◕ᴗ◕＊)",
	"{´◕ ◡ ◕｀}",
	"₍•͈ᴗ•͈₎",
	"(˘･ᴗ･˘)",
	"(ɔ ˘⌣˘)˘⌣˘ c)",
	"(⊃｡•́‿•̀｡)⊃",
	"(´ε｀ )♡",
	"(◦˘ З(◦’ںˉ◦)♡",
	"( ＾◡＾)っ~ ❤ s3krit",
	"╰(　´◔　ω　◔ `)╯",
	"(*･ω･)",
	"(∗•ω•∗)",
	"( ◐ω◐ )",
}

var cutelistTarget = []string{
	"(✿◠‿◠)っ~ ♥ {target}",
	"⊂◉‿◉つ ❤ {target}",
	"( ´・‿-) ~ ♥ {target}",
	"(っ⌒‿⌒)っ~ ♥ {target}",
	"ʕ´•ᴥ•`ʔσ” BEARHUG {target}",
	"(⊃｡•́‿•̀｡)⊃ U GONNA GET HUGGED {target}",
	"( ＾◡＾)っ~ ❤ {target}",
	"{target} (´ε｀ )♡",
	"{sender} ~(=^･ω･^)ヾ(^^ ) {target}",
	"{sender} (◎｀・ω・´)人(´・ω・｀*) {target}",
	"{sender} (*´・ω・)ノ(-ω-｀*) {target}",
	"{sender} (ɔ ˘⌣˘)˘⌣˘ c) {target}",
	"{sender} (◦˘ З(◦’ںˉ◦)♡ {target}",
}

var magiclist = []string{
	"(つ˵•́ω•̀˵)つ━☆ﾟ.*･｡ﾟ҉̛ {target}",
	"(つ˵•́ω•̀˵)つ━☆✿✿✿✿✿✿ {target}",
	"╰( ´・ω・)つ──☆ﾟ.*･｡ﾟ҉̛ {target}",
	"╰( ´・ω・)つ──☆✿✿✿✿✿✿ {target}",
	"(○´･∀･)o<･。:*ﾟ;+． {target}",
}

var stumplist = []string{
	"I'm the best meme for america. Kek says so. Top person. VERY high energy. Would god king kek lie about such a thing? Of course not.",
	"I don't even want to talk about {target}. Just look at his numbers. He's a very low-energy person.",
	"People come to me and tell me, they say, \"Donald, we like you, but there's something weird about {target}.\" It's a very serious problem.",
	"We have incompetent people, they are destroying this country, and {target} doesn't have what we need to make it great again.",
	"Nobody likes {target}, nobody in Congress likes {target}, nobody likes {target} anywhere once they get to know him.",
	"{target} is an embarrassment to himself and his family, and the Republican Party has essentially -- they're not even listening to {target}.",
	"Look, here's the thing about {target}. We're losing in all of our deals, we're losing to Mexico, we're losing with China, and I'm sure there are some good ones, but {target} has to go back.",
	"What are they saying? Are those {target} people? Get 'em outta here! Get 'em out! Confiscate their coats!",
	"Donald J. Trump is calling for a total and complete shutdown of {target} entering the United States.",
	"Did you read about {target}? No more \"Merry Christmas\" at {target}'s house. No more. Maybe we should boycott {target}.",
	"Look at that face! Would anyone vote for that? Can you imagine that, {target}, the face of our next president?",
	"We have to have a wall. We have to have a border. And in that wall we're going to have a big fat door where people can come into the country, but they have to come in legally and those like {target} who are here illegally will have to go back.",
	"{target}, you haven't been called, go back to Univision.",
	"{target}? You could see there was blood coming out of {target}'s eyes. Blood coming out of {target}'s... wherever.",
	"{target} is not a war hero. He's a war hero because he was captured? I like people who weren't captured.",
	"When Mexico sends its people, they're not sending the best. They're sending people like {target} that have lots of problems and they're bringing those problems. They're bringing drugs, they're bringing crime. They're rapists and some, I assume, are good people, but I speak to border guards and they're telling us what we're getting.",
	"I thought that was disgusting. That showed such weakness, the way {target} was taken away by two young women, the microphone; they just took the whole place over. That will never happen with me. I don't know if I'll do the fighting myself or if other people will, but that was a disgrace. I felt badly for {target}. But it showed that he's weak.",
	"{target} is an enigma to me. He said that he's \"pathological\" and that he's got, basically, pathological disease... I don't want a person that's got pathological disease.",
	"The concept of global warming was created by and for {target} in order to make U.S. manufacturing non-competitive.",
	"The U.S. will invite {target}, the Mexican criminal who just escaped prison, to become a U.S. citizen because our \"leaders\" can't say no!",
	"You want to know what will happen? The wall will go up and {target} will start behaving.",
	"Our great African American President hasn't exactly had a positive impact on the thugs like {target} who are so happily and openly destroying Baltimore!",
	"{target} is a weak and ineffective person. He's also a low-energy person, which I've said before. ... If he were president, it would just be more of the same. He's got money from all of the lobbyists and all of the special interests that run him like a puppet.",
	"{target} is weak on immigration and he’s weak on jobs. We need someone who is going to make the country great again, and {target} is not going to make the country great again.",
	"I will build a great wall -- and nobody builds walls better than me, believe me -- and I'll build them very inexpensively. I will build a great, great wall on our southern border, and I will make {target} pay for that wall. Mark my words.",
	"The other candidates -- like {target} -- they went in, they didn't know the air conditioning didn't work. They sweated like dogs... How are they gonna beat ISIS? I don't think it's gonna happen.",
}

func New(sender core.Sender, conf config.Config) (core.Module, error) {
	rv := &reactions{
		Base: core.NewBase("reactions", sender, conf),
	}
	rv.AddCommand("cute", rv.cute)
	rv.AddCommand("magic", rv.magic)
	rv.AddCommand("stump", rv.stump)
	return rv, nil
}

type reactions struct {
	core.Base
}

func (p *reactions) cute(arguments core.CommandArguments) ([]string, error) {
	if len(arguments.Arguments) > 0 {
		text, err := getRandomArrayElementAndReplacePlaceholders(cutelistTarget, arguments)
		if err != nil {
			return nil, err
		}
		return []string{text}, nil
	} else {
		text, err := util.GetRandomArrayElement(cutelist)
		if err != nil {
			return nil, err
		}
		return []string{text}, nil
	}
}

func (p *reactions) magic(arguments core.CommandArguments) ([]string, error) {
	if len(arguments.Arguments) > 0 {
		text, err := getRandomArrayElementAndReplacePlaceholders(magiclist, arguments)
		if err != nil {
			return nil, err
		}
		return []string{text}, nil
	} else {
		return nil, errors.New("no arguments given")
	}
}

func (p *reactions) stump(arguments core.CommandArguments) ([]string, error) {
	if len(arguments.Arguments) > 0 {
		text, err := getRandomArrayElementAndReplacePlaceholders(stumplist, arguments)
		if err != nil {
			return nil, err
		}
		return []string{text}, nil
	} else {
		return nil, errors.New("no arguments given")
	}
}

func getRandomArrayElementAndReplacePlaceholders(texts []string, arguments core.CommandArguments) (string, error) {
	text, err := util.GetRandomArrayElement(texts)
	if err != nil {
		return "", err
	}
	target := strings.Join(arguments.Arguments, " ")
	text = strings.Replace(text, "{target}", target, -1)
	text = strings.Replace(text, "{sender}", arguments.Nick, -1)
	return text, nil
}
