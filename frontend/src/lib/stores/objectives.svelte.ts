// simulate enum
export enum ObjectiveType {
  Hoisting = 0,
  Risk = 1,
}

interface IHoistingConfig {
}

interface IRiskConfig {
}

interface IObjectives {
  numberOfObjectives: ObjectiveType[];
  HoistingConfig?: IHoistingConfig;
  RiskConfig?: IRiskConfig;
}

interface IOptions {
  label: string;
  value: string;
}

class ObjectiveStore {
  objectives = $state<IObjectives>({
    numberOfObjectives: []
  })


  objectiveList: IOptions[] = [
    {
      label: 'Risk',
      value: 'risk'
    },
    {
      label: 'Hoisting',
      value: 'hoisting'
    },
    {
      label: 'Safety',
      value: 'safety'
    }
  ]
  selectObjectiveOptions: IOptions[] = []

  selectedObjectiveList: IOptions[] = []
  selectedObjectiveOptions: IOptions[] = []

  selectObjective = () => {

  }

  deselectObjective = () => {

  }

  addObjective(name: ObjectiveType) {
    switch (name) {
      case ObjectiveType.Hoisting:
        // Add Hoisting config

        this.objectives.numberOfObjectives.push(name);
        break

      case ObjectiveType.Risk:
        // Add Risk config

        this.objectives.numberOfObjectives.push(name);
        break

      default:
        return
    }

  }

}

export const objectiveStore = new ObjectiveStore();