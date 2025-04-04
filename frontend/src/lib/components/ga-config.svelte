<script lang="ts">

  import type {IGAConfig} from "$lib/stores/algorithms";

  interface Props {
    config: IGAConfig  | never
  }

  const {
    config,
  }: Props = $props()
</script>


<div>
  GA CONFIG
</div>