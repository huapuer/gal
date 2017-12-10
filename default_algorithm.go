package gal

import (
	"sort"
	"time"
	"math/rand"
)

type DefaultAlgorithm struct{
	PopulationSize int
	GoalScore float64
	MaxGen int64

	currentGeneration DefaultAlgorithms
	nextGeneration DefaultAlgorithms
	FitnessScore float64
	readIndex int
	writeIndex int
	TotalScore float64
}

func(da *DefaultAlgorithm) AppendCurrentGeneration(algorithm IAlgorithm){
	da.currentGeneration=append(da.currentGeneration, algorithm)
}

type DefaultAlgorithms []IAlgorithm

func (das DefaultAlgorithms) Len() int{
	return len(das)
}

func (das DefaultAlgorithms) Less(i, j int) bool{
	return das[i].GetFitnessScore() > das[j].GetFitnessScore()
}

func (das DefaultAlgorithms) Swap(i, j int){
	das[i], das[j] = das[j], das[i]
}

func (da *DefaultAlgorithm) Init(){
	rand.Seed(time.Now().UnixNano())

	da.currentGeneration=make([]IAlgorithm, 0, da.PopulationSize)
	da.nextGeneration=make([]IAlgorithm, da.PopulationSize, da.PopulationSize)
}

func (da *DefaultAlgorithm) NextIndividual() IAlgorithm{
	if da.readIndex == len(da.currentGeneration){
		da.readIndex=0
		return nil
	}
	ret := da.currentGeneration[da.readIndex]
	da.readIndex++
	return ret
}

func (da *DefaultAlgorithm) SortPopulation() {
	sort.Sort(da.currentGeneration)
}

func (da *DefaultAlgorithm) AchieveGoal(gen int64) IAlgorithm{
	if da.currentGeneration[0].GetFitnessScore() >= da.GoalScore || gen >= da.MaxGen{
		return da.currentGeneration[0]
	}
	return nil
}

func (da *DefaultAlgorithm) SelectIndividual() IAlgorithm{
	r:=rand.Float64()*da.TotalScore
	var roulette float64 = 0
	i:=0
	for _, a :=range da.currentGeneration{
		roulette+=a.GetFitnessScore()
		if roulette >= r{
			//fmt.Printf("SELECT: %d\n", i)
			return a.Duplicate()
		}
		i++
	}
	panic("Selecting Individual Failed.")
}

func (da *DefaultAlgorithm) FirstIndividual() IAlgorithm{
	return da.currentGeneration[0]
}

func (da *DefaultAlgorithm) ToGenerateNextPopulation() bool{
	return da.writeIndex < len(da.currentGeneration)
}

func (da *DefaultAlgorithm) PushIndividual(individual IAlgorithm){
	da.nextGeneration[da.writeIndex] = individual
	da.writeIndex++
}

func (da *DefaultAlgorithm) NextIterate(){
	da.readIndex=0
	da.writeIndex=0
	da.TotalScore = 0
	da.currentGeneration, da.nextGeneration = da.nextGeneration, da.currentGeneration
}

func(da *DefaultAlgorithm) GetFitnessScore() float64{
	return da.FitnessScore
}

func(da *DefaultAlgorithm) SumUpFitness(){
	for _, a :=range da.currentGeneration{
		da.TotalScore+=a.GetFitnessScore()
	}
}