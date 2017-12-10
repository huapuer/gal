package gal

import (
	"sync"
	"fmt"
)

type DefaultWorkflow struct{
	A IAlgorithm
}

func (dw DefaultWorkflow) Start(){
	dw.A.GeneratePopulation()
	gen:=0
	for {
		fmt.Printf("GEN: %d\n", gen)
		var wg sync.WaitGroup
		for {
			next:=dw.A.NextIndividual()
			if next!=nil {
				wg.Add(1)
				go func(){
					next.Act()
					wg.Done()
				}()
			}else{
				break
			}
		}
		wg.Wait()
		for {
			next:=dw.A.NextIndividual()
			if next!=nil {
				wg.Add(1)
				go func(){
					next.Fitness()
					wg.Done()
				}()
			}else{
				break
			}
		}
		wg.Wait()
		dw.A.SortPopulation()
		dw.A.SumUpFitness()
		best:=dw.A.AchieveGoal(int64(gen))
		if best!=nil{
			fmt.Printf("FINAL: %f\n", best.GetFitnessScore())
			//output result
			best.OutputResult()
			break;
		}
		fmt.Printf("BEST: %f\n", dw.A.FirstIndividual().GetFitnessScore())
		dw.A.PushIndividual(dw.A.FirstIndividual().(IAlgorithm).Duplicate().(IAlgorithm))
		for {
			i1:=dw.A.SelectIndividual().(IAlgorithm)
			i2:=dw.A.SelectIndividual().(IAlgorithm)
			p1:=i1.SelectChangePoint()
			p2:=i2.SelectChangePoint()
			if dw.A.ToGenerateNextPopulation() {
				if i1.Cross(p2) {
					i1.Mutate()
					dw.A.PushIndividual(i1)
				}
			}else{
				break
			}
			if dw.A.ToGenerateNextPopulation() {
				if i2.Cross(p1) {
					i2.Mutate()
					dw.A.PushIndividual(i2)
				}
			}else{
				break
			}
		}
		dw.A.NextIterate()
		gen++
	}
}