// simulate enum
export enum ObjectiveType {
  Hoisting = 0,
  Risk = 1,
  Safety  = 2,
}

interface IHoistingConfig {
}

interface IRiskConfig {
}

interface IObjectives {
  selectedObjectives: ObjectiveType[];
  HoistingConfig?: IHoistingConfig;
  RiskConfig?: IRiskConfig;
}

export interface IOptions {
  label: string;
  value: ObjectiveType;
  isChecked: boolean;
  content: string;
}

class ObjectiveStore {
  objectives = $state<IObjectives>({
    selectedObjectives: []
  })


  objectiveList= $state<IOptions[]>([
      {
        label: 'Risk',
        value: ObjectiveType.Risk,
        isChecked: false,
        content: "What is Risk Objective and How to calculate it?"
      },
      {
        label: 'Hoisting',
        value: ObjectiveType.Hoisting,
        isChecked: false,
        content: "What is Hoisting and How to calculate it?"
      },
      {
        label: 'Safety',
        value: ObjectiveType.Safety,
        isChecked: false,
        content: "What is Safety Objective and How to calculate it?"
      }
    ])


  selectObjectiveOption = $state<IOptions>()


  selectObjective = (option: IOptions) => {
    if (option.isChecked) {
      this.objectives.selectedObjectives.push(option.value)
    } else {
      this.objectives.selectedObjectives = this.objectives.selectedObjectives.filter(s => s !== option.value)
    }
  }

}

export const objectiveStore = new ObjectiveStore();