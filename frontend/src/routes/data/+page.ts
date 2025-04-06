import type { PageLoad } from './$types';
import {ProblemInfo} from "$lib/wailsjs/go/main/App";

export const load: PageLoad = async ({ params }) => {
  // get locations from problems

  const data = await ProblemInfo()

  console.log(data)

  return {

  };
};