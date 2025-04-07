<script lang="ts">

  import type {Snippet} from "svelte";

  interface Props {
    isModalOpen: boolean
    content: Snippet
    moreButtons: Snippet
    buttonText?: string
  }

  let {isModalOpen = $bindable(),
    content, moreButtons,
    buttonText }: Props = $props();

</script>

<div class="modal" class:modal-open={isModalOpen}>
  <div class="modal-box w-11/12 max-w-5xl">
    {@render content()}
    <div class="modal-action">
      <!-- 🔵 set false on click -->
      {@render moreButtons()}
      <button class="btn" onclick={()=>isModalOpen = false}>{buttonText ?? 'OK'}</button>
    </div>
  </div>
</div>