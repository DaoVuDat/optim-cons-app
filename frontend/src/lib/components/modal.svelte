<script lang="ts">

  import type {Snippet} from "svelte";

  interface Props {
    isModalOpen: boolean
    content: Snippet
    moreButtons?: Snippet
    buttonText?: string
    mainActionButton?: () => void
  }

  let {
    isModalOpen = $bindable(),
    content, moreButtons,
    buttonText, mainActionButton
  }: Props = $props();

</script>

<div class="modal" class:modal-open={isModalOpen}>
  <div class="modal-box w-11/12 max-w-5xl">
    {@render content()}
    <div class="modal-action">
      {#if moreButtons}
        {@render moreButtons()}
      {/if}
      <button class="btn" onclick={()=>{
        mainActionButton?.()
        isModalOpen = false
      }}>{buttonText ?? 'OK'}</button>
    </div>
  </div>
</div>