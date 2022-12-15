package hw03frequencyanalysis

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Change to true if needed.
var taskWithAsteriskIsCompleted = true
var text = `Как видите, он  спускается  по  лестнице  вслед  за  своим
	другом   Кристофером   Робином,   головой   вниз,  пересчитывая
	ступеньки собственным затылком:  бум-бум-бум.  Другого  способа
	сходить  с  лестницы  он  пока  не  знает.  Иногда ему, правда,
		кажется, что можно бы найти какой-то другой способ, если бы  он
	только   мог   на  минутку  перестать  бумкать  и  как  следует
	сосредоточиться. Но увы - сосредоточиться-то ему и некогда.
		Как бы то ни было, вот он уже спустился  и  готов  с  вами
	познакомиться.
	- Винни-Пух. Очень приятно!
		Вас,  вероятно,  удивляет, почему его так странно зовут, а
	если вы знаете английский, то вы удивитесь еще больше.
		Это необыкновенное имя подарил ему Кристофер  Робин.  Надо
	вам  сказать,  что  когда-то Кристофер Робин был знаком с одним
	лебедем на пруду, которого он звал Пухом. Для лебедя  это  было
	очень   подходящее  имя,  потому  что  если  ты  зовешь  лебедя
	громко: "Пу-ух! Пу-ух!"- а он  не  откликается,  то  ты  всегда
	можешь  сделать вид, что ты просто понарошку стрелял; а если ты
	звал его тихо, то все подумают, что ты  просто  подул  себе  на
	нос.  Лебедь  потом  куда-то делся, а имя осталось, и Кристофер
	Робин решил отдать его своему медвежонку, чтобы оно не  пропало
	зря.
		А  Винни - так звали самую лучшую, самую добрую медведицу
	в  зоологическом  саду,  которую  очень-очень  любил  Кристофер
	Робин.  А  она  очень-очень  любила  его. Ее ли назвали Винни в
	честь Пуха, или Пуха назвали в ее честь - теперь уже никто  не
	знает,  даже папа Кристофера Робина. Когда-то он знал, а теперь
	забыл.
		Словом, теперь мишку зовут Винни-Пух, и вы знаете почему.
		Иногда Винни-Пух любит вечерком во что-нибудь поиграть,  а
	иногда,  особенно  когда  папа  дома,  он больше любит тихонько
	посидеть у огня и послушать какую-нибудь интересную сказку.
		В этот вечер...`
var text2 = `Жил старик со своею старухой
	У самого синего моря;
	Они жили в ветхой землянке
	Ровно тридцать лет и три года.
	Старик ловил неводом рыбу,
	Старуха пряла свою пряжу.
	Раз он в море закинул невод, -
	Пришел невод с одною тиной.
	Он в другой раз закинул невод, -
	Пришел невод с травой морскою.
	В третий раз закинул он невод, -
	Пришел невод с одною рыбкой,
	С непростою рыбкой, - золотою.
	Как взмолится золотая рыбка!
	Голосом молвит человечьим:
	"Отпусти ты, старче, меня в море,
	Дорогой за себя дам откуп:
	Откуплюсь чем только пожелаешь".
	Удивился старик, испугался:
	Он рыбачил тридцать лет и три года
	И не слыхивал, чтоб рыба говорила.
	Отпустил он рыбку золотую
	И сказал ей ласковое слово:
	"Бог с тобою, золотая рыбка!
	Твоего мне откупа не надо;
	Ступай себе в синее море,
	Гуляй там себе на просторе".`
var text3 = `cat and dog, one and one man`

func TestTop10(t *testing.T) {
	TaskWithAsteriskIsCompleted = taskWithAsteriskIsCompleted
	t.Run("no words in empty string", func(t *testing.T) {
		require.Len(t, Top10(""), 0)
	})
	if taskWithAsteriskIsCompleted {
		t.Run("no words in empty string", func(t *testing.T) {
			require.Len(t, Top10(";-."), 0)
		})
	}
	t.Run("positive test", func(t *testing.T) {
		if taskWithAsteriskIsCompleted {
			expected := []string{
				"а",         // 8
				"он",        // 8
				"и",         // 6
				"ты",        // 5
				"что",       // 5
				"в",         // 4
				"его",       // 4
				"если",      // 4
				"кристофер", // 4
				"не",        // 4
			}
			require.Equal(t, expected, Top10(text))
		} else {
			expected := []string{
				"он",        // 8
				"а",         // 6
				"и",         // 6
				"ты",        // 5
				"что",       // 5
				"-",         // 4
				"Кристофер", // 4
				"если",      // 4
				"не",        // 4
				"то",        // 4
			}
			require.Equal(t, expected, Top10(text))
		}
	})
	t.Run("positive test", func(t *testing.T) {
		if taskWithAsteriskIsCompleted {
			expected := []string{
				"в",
				"невод",
				"он",
				"с",
				"и",
				"закинул",
				"море",
				"пришел",
				"раз",
				"старик",
			}
			require.Equal(t, expected, Top10(text2))
		} else {
			expected := []string{
				"в",
				"-",
				"с",
				"Пришел",
				"закинул",
				"невод",
				"невод,",
				"он",
				"И",
				"Он",
			}
			require.Equal(t, expected, Top10(text2))
		}
	})
	t.Run("positive test", func(t *testing.T) {
		if taskWithAsteriskIsCompleted {
			expected := []string{
				"and",
				"one",
				"cat",
				"dog",
				"man",
			}
			require.Equal(t, expected, Top10(text3))
		} else {
			expected := []string{
				"and",
				"one",
				"cat",
				"dog,",
				"man",
			}
			require.Equal(t, expected, Top10(text3))
		}
	})
}
