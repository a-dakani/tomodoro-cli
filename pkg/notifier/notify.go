package notifier

import (
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/gen2brain/beeep"
	"os"
	"time"
)

type Notifier struct {
	Title         string
	ImagePath     string
	AudioFilePath string
}

var notifier *Notifier

func NewNotifier(title, imagePath, audioFilePath string) *Notifier {
	notifier = &Notifier{
		Title:         title,
		ImagePath:     imagePath,
		AudioFilePath: audioFilePath,
	}

	return notifier
}

func (n *Notifier) beep() error {
	f, err := os.Open(n.AudioFilePath)
	if err != nil {
		return err
	}

	// Decode the audio file
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		return err
	}

	// Initialize the speaker
	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		return err
	}

	// Play the audio
	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	// Wait for the audio to finish playing
	<-done

	// Close the audio file and the speaker
	err = streamer.Close()
	if err != nil {
		return err
	}

	return nil
}

func (n *Notifier) systemNotify(message string) error {
	err := beeep.Notify(n.Title, message, n.ImagePath)
	if err != nil {
		panic(err)
	}
	return nil
}

func (n *Notifier) Notify(message string) {
	// run in background

	go func() {
		_ = n.systemNotify(message)
		//_ = n.beep()
	}()

}
