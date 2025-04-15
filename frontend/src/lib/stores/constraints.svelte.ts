import {data} from "$lib/wailsjs/go/models";
import {
  coverInCraneRadiusConfig,
  type ICoverInCraneRadiusConfig,
  type IInclusiveZoneConfig, inclusiveZoneConfig,
  type IOutOfBoundConfig,
  type IOverlapConfig, outOfBoundConfig, overlapConfig
} from "$lib/stores/constraints";
import {objectiveStore} from "$lib/stores/objectives.svelte";

type IConfigType = IOutOfBoundConfig | IOverlapConfig | ICoverInCraneRadiusConfig | IInclusiveZoneConfig

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
    }
  ])

  validConstraintList = $derived.by<IConstraintOptions[]>(() => {
    // filtering the "cover in crane radius"
    const selectedObjectives = objectiveStore.objectives.selectedObjectives.map(selectedObjective => selectedObjective.objectiveType)
    return this.constraintList.filter(cons => cons.value !== data.ConstraintType.CoverInCraneRadius ||
      selectedObjectives.includes(data.ObjectiveType.HoistingObjective))
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
    }
  }

}

export const constraintsStore = new ConstraintsStore();