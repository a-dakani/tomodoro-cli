package notifier

import (
	"github.com/gen2brain/beeep"
)

type Notifier struct {
	Title     string
	ImagePath string
}

var notifier *Notifier

func NewNotifier(title, imagePath string) *Notifier {
	notifier = &Notifier{
		Title:     title,
		ImagePath: imagePath,
	}

	return notifier
}

func (n *Notifier) Notify(message string) error {
	err := beeep.Alert(n.Title, message, n.ImagePath)
	if err != nil {
		return err
	}
	return nil
}
