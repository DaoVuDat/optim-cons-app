

export function roundNDecimal (value: number, n: number) {
    const pow = Math.pow(10, n)
    return Math.round((value + Number.EPSILON) * pow) / pow
}