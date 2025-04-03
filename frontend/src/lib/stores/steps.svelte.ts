import {objectiveStore} from "$lib/stores/objectives.svelte";

class StepStore {
  step = $state<number>(1)

  stepsList = [
    {
      number: 1,
      name: 'Objectives'
    },
    {
      number: 2,
      name: 'Algorithm'
    },
    {
      number: 3,
      name: 'Problem'
    },
    {
      number: 4,
      name: 'Optimize'
    }
  ]

  nextStep = () => {
    if (this.step >= 4) {
      return
    }
    this.step++

  }

  prevStep = () => {
    if (this.step <= 1) {
      return
    }
    this.step--
  }

}

export const stepStore = new StepStore()