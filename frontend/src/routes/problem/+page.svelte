<script lang="ts">
    import clsx from "clsx";
    import {stepStore} from "$lib/stores/steps.svelte.js";
    import {problemList, problemStore, type ProblemWithLabel} from "$lib/stores/problem.svelte.js";
    import continuousProblemConfigComponent from "$lib/components/problem-configs/continuous-config.svelte";
    import gridProblemConfigComponent from "$lib/components/problem-configs/grid-config.svelte";
    import PredeterminatedConfig from "$lib/components/problem-configs/predeterminated-config.svelte";
    import {goto} from "$app/navigation";
    import {CreateProblem} from "$lib/wailsjs/go/main/App";
    import {objectives} from "$lib/wailsjs/go/models";

    const configComponents = {
        [objectives.ProblemType.ContinuousConstructionLayout]: continuousProblemConfigComponent,
        [objectives.ProblemType.GridConstructionLayout]: gridProblemConfigComponent,
        [objectives.ProblemType.PredeterminedConstructionLayout]: PredeterminatedConfig,
    }

    const component = $derived.by(() => {
        if (problemStore.getValidSelection()) {
            return configComponents[problemStore.selectedProblem!.value]
        }
    })

    let loading = $state<boolean>(false)

    const handleClick = (prob: ProblemWithLabel) => {
        problemStore.selectedProblem = prob;
    }

    const handleNext = async () => {
        loading = true
        // TODO: add GRID problem and PREDETERMINATED LOCATIONS problem
        if (problemStore.selectedProblem) {
            switch (problemStore.selectedProblem!.value) {
                case objectives.ProblemType.ContinuousConstructionLayout :
                    const config = problemStore.getConfig(problemStore.selectedProblem!.value)

                    await CreateProblem(problemStore.selectedProblem!.value,
                        config.length,
                        config.width,
                        config.facilitiesFilePath.value,
                        config.phasesFilePath.value)
                    break
                case objectives.ProblemType.GridConstructionLayout :
                    break
                case objectives.ProblemType.PredeterminedConstructionLayout :
                    break
            }
        }

        loading = false
        await goto('/data')
        stepStore.nextStep()
    }


</script>

<div class="h-[calc(100vh-64px-64px)] w-full text-lg pt-4 flex flex-col justify-between items-center">
    <!-- Top Section -->
    <section class="mt-8 text-black">
        <h1 class="text-5xl font-bold">Select problem</h1>
    </section>

    <!-- Content -->
    <section class="px-24 grid grid-cols-12 gap-4 w-[1400px] auto-rows-min">
        <div
                class="h-[420px] bg-base-100 px-2 py-4 card shadow-md rounded-lg col-span-4 flex flex-col space-y-2 overflow-y-auto">
            {#each problemList as prob (prob)}
                <button class={clsx("p-4 rounded h-12 flex justify-between items-center cursor-pointer text-left",
        problemStore.selectedProblem?.value === prob.value ? 'bg-[#422AD5] text-white' : '')}
                        onclick={() => handleClick(prob)}>
                    {prob.label}
                </button>
            {/each}
        </div>
        <div
                class="h-[420px] bg-base-100 overflow-y-auto card p-4 shadow-md rounded-lg col-span-8 flex flex-col justify-center items-center">
            {#if problemStore.getValidSelection()}
                {@const Component = component}
                <Component/>
            {:else}
                <p>Please select problem</p>
            {/if}
        </div>

        {#if loading}
            <div class="toast toast-center toast-middle">
                <div class="alert alert-info">
                    <span>Loading data...</span>
                </div>
            </div>
        {/if}
    </section>

    <!-- Bottom Section -->
    <section class="w-full text-end">
        <a class="ml-4 btn" href="/" onclick={() => stepStore.prevStep()}>Back</a>
        <button class={clsx('ml-4 btn', problemStore.getValidSelection() ? '' : 'btn-disabled')}
                onclick={() => handleNext()}
        >Next
        </button>
    </section>
</div>