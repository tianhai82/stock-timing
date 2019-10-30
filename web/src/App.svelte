<script>
  import { Select, Button } from "smelte";
  import Signin from "./components/Signin.svelte";

  export let name;
  let showSignIn = false;
  let stock;
  let instruments = ["Apple", "Amazon"];
  fetch("https://stock-timing.web.app/rpc/instruments")
    .then(resp => resp.json())
    .then(message => {
      console.log(message);
    });
</script>

<div class="container mx-auto h-full items-center">
  <header class="bg-teal-200 p-3 shadow">
    <div class="flex flex-row align-middle">
      <img
        src="/images/time-money.png"
        alt="logo"
        class="object-contain h-10 mx-4 mt-1" />
      <h4>Stock Timing</h4>
      <Button on:click={() => (showSignIn = true)}>Sign In</Button>
    </div>
  </header>
  <Signin bind:showSignIn />
  <div class="p-4">
    <Select
      bind:value={stock}
      outlined
      autocomplete
      label="Enter Company Name"
      items={instruments} />
    {stock}
  </div>
</div>
