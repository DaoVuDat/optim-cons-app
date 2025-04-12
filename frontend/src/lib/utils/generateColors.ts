export function generateLabelFriendlyColors( n: number, darkLabels = true) {
    // Pre-defined color palette with good contrast for labels
    // These colors are selected to work well with white text (if darkLabels is false)
    // or dark text (if darkLabels is true)
    const baseColors = darkLabels ?
        // Lighter colors for dark labels
        ['#8dd3c7', '#ffffb3', '#bebada', '#fb8072', '#80b1d3', '#fdb462', '#b3de69',
            '#fccde5', '#d9d9d9', '#bc80bd', '#ccebc5', '#ffed6f', '#a6cee3', '#b2df8a'] :
        // Darker colors for white labels
        ['#1f77b4', '#ff7f0e', '#2ca02c', '#d62728', '#9467bd', '#8c564b', '#e377c2',
            '#7f7f7f', '#bcbd22', '#17becf', '#756bb1', '#636363', '#31a354', '#ad494a'];

    const result = [];

    // If we need more colors than in our base palette
    if (n <= baseColors.length) {
        console.log("run that")
        return baseColors.slice(0, n);
    } else {
        console.log("run this")
        // Include all base colors
        result.push(...baseColors);

        // Generate additional colors with slight variations
        for (let i = baseColors.length; i < n; i++) {
            const baseColor = baseColors[i % baseColors.length];

            // Parse the hex color to RGB
            const r = parseInt(baseColor.slice(1, 3), 16);
            const g = parseInt(baseColor.slice(3, 5), 16);
            const b = parseInt(baseColor.slice(5, 7), 16);

            // Create variation - shift hue slightly based on how many extra colors we've made
            const multiplier = Math.floor(i / baseColors.length);

            // Adjust the color components based on whether we need lighter or darker variations
            let newR = (r + 20 * multiplier) % 256;
            let newG = (g + 30 * multiplier) % 256;
            let newB = (b + 40 * multiplier) % 256;

            // For dark labels, ensure the color is light enough
            if (darkLabels) {
                newR = Math.max(newR, 160);
                newG = Math.max(newG, 160);
                newB = Math.max(newB, 160);
            }
            // For white labels, ensure the color is dark enough
            else {
                newR = Math.min(newR, 200);
                newG = Math.min(newG, 200);
                newB = Math.min(newB, 200);
            }

            // Convert back to hex
            const newColor = '#' +
                newR.toString(16).padStart(2, '0') +
                newG.toString(16).padStart(2, '0') +
                newB.toString(16).padStart(2, '0');

            result.push(newColor);
        }

        return result;
    }

    // const colors: string[] = [];
    //
    // for (let i = 0; i < n; i++) {
    //     // Generate random RGB values within dark range
    //     // Keep values low to ensure dark colors (0-130)
    //     const r = Math.floor(Math.random() * 130);
    //     const g = Math.floor(Math.random() * 130);
    //     const b = Math.floor(Math.random() * 130);
    //
    //     // Convert to hex and ensure 2 digits with padStart
    //     const hexR = r.toString(16).padStart(2, '0');
    //     const hexG = g.toString(16).padStart(2, '0');
    //     const hexB = b.toString(16).padStart(2, '0');
    //
    //     // Add hex color string to array
    //     colors.push(`#${hexR}${hexG}${hexB}`);
    // }
    //
    // return colors;
}