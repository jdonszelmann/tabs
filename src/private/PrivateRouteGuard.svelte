<script>
    import { useNavigate, useLocation } from "svelte-navigator";
    import { user } from "../typescript/User";

    const navigate = useNavigate();
    const location = useLocation();

    $: if ($user === null) {
        navigate("/login", {
            state: { from: $location.pathname },
            replace: true,
        });
    }
</script>

{#if $user}
    <slot />
{/if}
