package agent

type Agent struct {
	Classifier Classifier
	Router     Router
}

func (a *Agent) Handle(input string) (string, error){
	intent, err := a.Classifier.Classify(input)
	if err != nil {
		return "", err
	}

	return a.Router.Route(intent, input)
}