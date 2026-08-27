<template>
  <div class="min-h-screen bg-slate-50 text-slate-900">
    <header class="border-b border-slate-200 bg-white shadow-xs">
      <div class="mx-auto flex max-w-5xl flex-wrap items-center justify-between gap-4 px-4 py-3 sm:px-6">
        <div class="flex items-center gap-2.5">
          <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-blue-600 font-bold text-white shadow-xs text-xs">
            WRS
          </div>
          <div>
            <span class="block text-sm font-bold text-slate-800 leading-tight">Welfare Registration</span>
            <span class="block text-[10px] text-slate-400">ระบบลงทะเบียนสวัสดิการแห่งรัฐ</span>
          </div>
        </div>

        <nav class="flex gap-1 rounded-lg border border-slate-200 bg-slate-100/80 p-1 text-xs sm:text-sm shadow-2xs">
          <RouterLink
            to="/register"
            active-class="bg-white text-blue-600 shadow-2xs font-semibold"
            class="rounded-md px-3 py-1.5 text-slate-600 transition hover:text-slate-900"
          >
            ลงทะเบียน (Register)
          </RouterLink>
          <RouterLink
            to="/status"
            active-class="bg-white text-blue-600 shadow-2xs font-semibold"
            class="rounded-md px-3 py-1.5 text-slate-600 transition hover:text-slate-900"
          >
            ติดตามสถานะ (Status)
          </RouterLink>
          <RouterLink
            to="/officer"
            active-class="bg-white text-blue-600 shadow-2xs font-semibold"
            class="rounded-md px-3 py-1.5 text-slate-600 transition hover:text-slate-900"
          >
            เจ้าหน้าที่ (Officer)
          </RouterLink>
        </nav>
      </div>
    </header>

    <div class="mx-auto max-w-5xl">
      <RouterView />
    </div>
  </div>
</template>
