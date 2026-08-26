<template>
  <header class="mx-auto mb-8 flex max-w-5xl flex-wrap items-center justify-between gap-4">
    <span class="font-semibold text-slate-900">Welfare Registration System</span>

    <nav class="flex gap-1 rounded-lg border bg-white p-1 text-sm shadow-sm">
      <RouterLink
        to="/register"
        active-class="bg-blue-600 text-white"
        class="rounded-md px-3 py-2 text-slate-600 transition hover:bg-slate-100"
      >
        Register
      </RouterLink>
      <RouterLink
        to="/status"
        active-class="bg-blue-600 text-white"
        class="rounded-md px-3 py-2 text-slate-600 transition hover:bg-slate-100"
      >
        Status
      </RouterLink>
      <RouterLink
        to="/officer"
        active-class="bg-blue-600 text-white"
        class="rounded-md px-3 py-2 text-slate-600 transition hover:bg-slate-100"
      >
        Officer
      </RouterLink>
    </nav>
  </header>

  <RouterView />
</template>
