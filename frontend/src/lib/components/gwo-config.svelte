<script lang="ts">

  import type {IGWOConfig} from "$lib/stores/algorithms";

  interface Props {
    config: IGWOConfig | never
  }

  const {
    config,
  }: Props = $props()
</script>


<div>
  IGWO CONFIG
</div>