package helper

import "math/rand"

func RandomNickname() string {

	arrayNick := [16]string{"dog", "cat", "butterfly", "duck", "cow", "monkey", "snake", "ant", "bear", "dolphin","bat","bird","crap","chicken","deer","eagle"}
	rnNumber := rand.Intn(16)
	rnNick := arrayNick[rnNumber]

	return rnNick
}

func RandomColor() string {

	arrayColor := [10]string{"#ff0000", "#ff4000", "#ffbf00", "#00ff00", "#00ffff", "#0040ff", "#8000ff", "#ff00ff", "#ff0040", "#808080"}
	rnNumber2 := rand.Intn(10)
	rnColor := arrayColor[rnNumber2]

	return rnColor
}
