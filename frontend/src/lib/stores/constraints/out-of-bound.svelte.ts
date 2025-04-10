

export interface IOutOfBoundConfig {
    AlphaOutOfBoundaryPenalty: number
    PowerDifferencePenalty: number
}


export const outOfBoundConfig = $state<IOutOfBoundConfig>({
    AlphaOutOfBoundaryPenalty: 20000,
    PowerDifferencePenalty: 1
})