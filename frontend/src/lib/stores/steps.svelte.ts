
class StepStore {
  step = $state<number>(1)

  stepsList = [
    {
      number: 1,
      name: 'Objectives'
    },
    {
      number: 2,
      name: 'Problem'
    },
    {
      number: 3,
      name: 'Configuration'
    },
    {
      number: 4,
      name: 'Constraints'
    },
    {
      number: 5,
      name: 'Algorithm'
    },
    {
      number: 6,
      name: 'Optimize'
    }
  ]

  nextStep = () => {
    if (this.step >= this.stepsList.length) {
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