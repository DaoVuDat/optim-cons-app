<script lang="ts">
  import "../app.css"
  import {stepStore} from "$lib/stores/steps.svelte";

  import clsx from "clsx";
  import {onMount} from "svelte";
  import {browser} from "$app/environment";
  import {SvelteToast} from "@zerodevx/svelte-toast";

  const {children} = $props()

  const options = {
    duration: 2000,
  }

  onMount(() => {
      if (browser) {

          const preventContextMenu = (e) => e.preventDefault();
          window.addEventListener('contextmenu', preventContextMenu);

          return () => {
              window.removeEventListener('contextmenu', preventContextMenu);
          };
      }
  });
</script>

<SvelteToast {options} />
<main class="p-8 bg-gray-100">
  <!-- Steps -->
  <div class="w-full flex justify-center items-centers text-black">
    <ul class="steps">
      {#each stepStore.stepsList as s}
          <li class={clsx('step w-48', s.number <= stepStore.step ? "step-primary": "")}>{s.name}</li>
      {/each}
    </ul>
  </div>
  <!-- Page -->
  {@render children()}
</main>
