package gal

type IAlgorithm interface{
	Init()
	NextIndividual() IAlgorithm
	SortPopulation()
	AchieveGoal(gen int64) IAlgorithm
	SelectIndividual() IAlgorithm
	FirstIndividual() IAlgorithm
	ToGenerateNextPopulation() bool
	PushIndividual(individual IAlgorithm)
	NextIterate()
	GetFitnessScore() float64

	GeneratePopulation()
	SelectChangePoint() interface{}
	Duplicate() IAlgorithm
	Cross(point interface{}) bool
	Mutate() bool
	Act()
	Fitness()
	SumUpFitness()
	OutputResult()
}