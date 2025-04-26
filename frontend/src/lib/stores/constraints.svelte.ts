import {data} from "$lib/wailsjs/go/models";
import {
  coverInCraneRadiusConfig,
  type ICoverInCraneRadiusConfig,
  type IInclusiveZoneConfig, inclusiveZoneConfig,
  type IOutOfBoundConfig,
  type IOverlapConfig, outOfBoundConfig, overlapConfig
} from "$lib/stores/constraints";
import {objectiveStore} from "$lib/stores/objectives.svelte";
import {type ISizeConfig, sizeConfig} from "$lib/stores/constraints/size.svelte";
import {problemStore} from "$lib/stores/problem.svelte";

type IConfigType = IOutOfBoundConfig | IOverlapConfig | ICoverInCraneRadiusConfig | IInclusiveZoneConfig | ISizeConfig

interface IConstraint {
  selectedConstraints: {
    constraintType: data.ConstraintType,
    config?: IConfigType
  }[];
}

export interface IConstraintOptions {
  label: string;
  value: data.ConstraintType;
  isChecked: boolean;
}

export type ConstraintConfigMap = {
  [data.ConstraintType.CoverInCraneRadius]: ICoverInCraneRadiusConfig;
  [data.ConstraintType.InclusiveZone]: IInclusiveZoneConfig;
  [data.ConstraintType.Overlap]: IOverlapConfig;
  [data.ConstraintType.OutOfBound]: IOutOfBoundConfig;
  [data.ConstraintType.Size]: ISizeConfig;
}

class ConstraintsStore {
  constraints = $state<IConstraint>({
    selectedConstraints: []
  })


  constraintList = $state<IConstraintOptions[]>([
    {
      label: 'Out of boundary',
      value: data.ConstraintType.OutOfBound,
      isChecked: false,
    },
    {
      label: 'Overlap',
      value: data.ConstraintType.Overlap,
      isChecked: false,
    },
    {
      label: 'Cover in crane radius',
      value: data.ConstraintType.CoverInCraneRadius,
      isChecked: false,
    },
    {
      label: 'Inclusive zone',
      value: data.ConstraintType.InclusiveZone,
      isChecked: false,
    },
    {
      label: 'Size',
      value: data.ConstraintType.Size,
      isChecked: false,
    }
  ])

  validConstraintList = $derived.by<IConstraintOptions[]>(() => {
    if (problemStore.selectedProblem?.value === data.ProblemName.PredeterminedConstructionLayout) {
      let constraints = []
      constraints.push(this.constraintList.find(cons => cons.value === data.ConstraintType.Size)!)

      return constraints
    } else {
      // filtering the "cover in crane radius"
      const selectedObjectives = objectiveStore.objectives.selectedObjectives.map(selectedObjective => selectedObjective.objectiveType)
      return this.constraintList
          .filter(cons => cons.value !== data.ConstraintType.Size) // fitering the Size constraint ( for predetermined problem)
          
    }


  })

  clearConstraint = () => {
    this.constraints.selectedConstraints.length = 0
    this.validConstraintList.forEach(constraint => {
      constraint.isChecked = false
    })
  }

  selectConstraint = (option: IConstraintOptions) => {
    if (option.isChecked) {
      const config = this.getConfig(option.value)

      this.constraints.selectedConstraints.push({
        constraintType: option.value,
        config
      })
    } else {
      this.constraints.selectedConstraints = this.constraints.selectedConstraints.filter(s => s.constraintType !== option.value)
    }
  }


  getConfig = <T extends data.ConstraintType>(type: T): ConstraintConfigMap[T] => {
    switch (type) {
      case data.ConstraintType.OutOfBound:
        return outOfBoundConfig as ConstraintConfigMap[T]
      case data.ConstraintType.Overlap:
        return overlapConfig as ConstraintConfigMap[T]
      case data.ConstraintType.InclusiveZone:
        return inclusiveZoneConfig as ConstraintConfigMap[T]
      case data.ConstraintType.CoverInCraneRadius:
        return coverInCraneRadiusConfig as ConstraintConfigMap[T]
      case data.ConstraintType.Size:
        return sizeConfig as ConstraintConfigMap[T]
    }
  }

}

export const constraintsStore = new ConstraintsStore();