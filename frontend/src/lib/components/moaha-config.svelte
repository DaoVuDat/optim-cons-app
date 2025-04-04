<script lang="ts">

  import type {IMOAHAConfig} from "$lib/stores/algorithms";

  interface Props {
    config: IMOAHAConfig | never
  }

  const {
    config,
  }: Props = $props()
</script>


<div>
  MOAHA CONFIG
</div>