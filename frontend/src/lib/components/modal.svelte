<script lang="ts">

  import type {Snippet} from "svelte";

  interface Props {
    isModalOpen: boolean
    content: Snippet
  }

  let {isModalOpen = $bindable(), content }: Props = $props();

</script>

<div class="modal" class:modal-open={isModalOpen}>
  <div class="modal-box w-11/12 max-w-5xl">
    {@render content()}
    <div class="modal-action">
      <!-- 🔵 set false on click -->
      <button class="btn" onclick={()=>isModalOpen = false}>OK</button>
    </div>
  </div>
</div>