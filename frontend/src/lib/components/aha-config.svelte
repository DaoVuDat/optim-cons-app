<script lang="ts">

  import type {IAHAConfig} from "$lib/stores/algorithms";

  interface Props {
    config: IAHAConfig | never
  }

  const {
    config,
  }: Props = $props()
</script>


<div>
  AHA CONFIG
</div>