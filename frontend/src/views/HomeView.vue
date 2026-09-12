<template>
  <!-- Custom Home Content: Full Page Mode -->
  <div v-if="hasHomeContent" class="min-h-screen">
    <!-- iframe mode -->
    <iframe
      v-if="isHomeContentUrl"
      :src="homeContent.trim()"
      class="h-screen w-full border-0"
      allowfullscreen
    ></iframe>
    <!-- HTML mode - SECURITY: homeContent is admin-only setting, XSS risk is acceptable -->
    <div v-else v-html="homeContent"></div>
  </div>

  <!-- Compact Home Page -->
  <div
    v-else-if="compactHomeEnabled"
    data-testid="compact-home"
    :style="brandStyleVars"
    class="relative flex min-h-screen flex-col overflow-x-clip brand-wash text-gray-900 dark:text-white"
  >
    <!-- Ambient background -->
    <div class="pointer-events-none absolute inset-0 overflow-hidden" aria-hidden="true">
      <div
        class="absolute left-1/2 top-0 h-[30rem] w-[30rem] -translate-x-1/2 -translate-y-1/3 rounded-full brand-orb-from blur-3xl"
      ></div>
      <div
        class="absolute bottom-0 right-0 h-72 w-72 translate-x-1/4 translate-y-1/4 rounded-full brand-orb-to blur-3xl"
      ></div>
      <div class="hero-grid absolute inset-0"></div>
    </div>

    <header class="relative z-10 px-4 py-4 sm:px-6">
      <nav class="mx-auto flex max-w-5xl flex-wrap items-center justify-between gap-3 sm:gap-4">
        <div class="flex min-w-0 flex-1 items-center gap-2.5">
          <img
            :src="siteLogo || '/logo.svg'"
            alt="Logo"
            class="h-9 w-9 shrink-0 rounded-xl object-contain shadow-sm"
          />
          <span class="min-w-0 truncate text-base font-semibold tracking-tight">{{ siteName }}</span>
        </div>
        <div class="flex max-w-full shrink-0 flex-wrap items-center justify-end gap-1.5">
          <LocaleSwitcher />
          <a
            v-if="docUrl"
            :href="docUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-900 dark:text-dark-400 dark:hover:bg-white/5 dark:hover:text-white"
            :title="t('home.viewDocs')"
          >
            <Icon name="book" size="md" />
          </a>
          <router-link
            v-if="showModelPlazaEntry"
            to="/model-plaza"
            class="flex h-9 shrink-0 items-center gap-1.5 rounded-lg px-2.5 text-sm font-medium text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-900 dark:text-dark-400 dark:hover:bg-white/5 dark:hover:text-white"
            :title="t('nav.modelPlaza')"
          >
            <Icon name="grid" size="md" />
            <span class="hidden sm:inline">{{ t('nav.modelPlaza') }}</span>
          </router-link>
          <button
            class="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-900 dark:text-dark-400 dark:hover:bg-white/5 dark:hover:text-white"
            :title="isDark ? t('home.switchToLight') : t('home.switchToDark')"
            :aria-label="isDark ? t('home.switchToLight') : t('home.switchToDark')"
            @click="toggleTheme"
          >
            <Icon v-if="isDark" name="sun" size="md" />
            <Icon v-else name="moon" size="md" />
          </button>
          <router-link
            :to="isAuthenticated ? dashboardPath : '/login'"
            class="ml-1 shrink-0 rounded-full brand-grad brand-grad-hover brand-shadow px-4 py-2 text-sm font-medium text-white transition-all duration-200 active:scale-[0.98]"
          >
            {{ isAuthenticated ? t('home.dashboard') : t('home.login') }}
          </router-link>
        </div>
      </nav>
    </header>

    <main class="relative z-10 flex min-w-0 flex-1 items-center justify-center px-4 py-14 sm:px-6">
      <div class="w-full min-w-0 max-w-lg text-center">
        <div class="relative mx-auto mb-7 h-24 w-24">
          <span
            class="absolute inset-0 rounded-3xl brand-grad opacity-40 blur-xl"
          ></span>
          <img
            :src="siteLogo || '/logo.svg'"
            alt="Logo"
            class="relative h-24 w-24 rounded-3xl object-contain shadow-xl ring-1 ring-black/5 dark:ring-white/10"
          />
        </div>

        <h1 class="[overflow-wrap:anywhere] text-3xl font-bold tracking-tight sm:text-4xl">
          {{ siteName }}
        </h1>
        <p
          class="mx-auto mt-4 max-w-md whitespace-pre-wrap text-pretty [overflow-wrap:anywhere] text-base leading-relaxed text-gray-600 dark:text-dark-300"
        >
          {{ siteSubtitle }}
        </p>

        <div class="mt-8 flex flex-col items-center justify-center">
          <router-link
            :to="isAuthenticated ? dashboardPath : '/login'"
            class="inline-flex w-full items-center justify-center gap-2 rounded-full brand-grad brand-grad-hover brand-shadow px-8 py-3 text-base font-medium text-white transition-all duration-200 active:scale-[0.98] sm:w-auto"
          >
            {{ isAuthenticated ? t('home.goToDashboard') : t('home.getStarted') }}
            <Icon name="arrowRight" size="md" :stroke-width="2" />
          </router-link>
        </div>

        <div class="mt-8 flex flex-wrap items-center justify-center gap-x-5 gap-y-2.5">
          <span
            v-for="tag in heroTags"
            :key="tag.label"
            class="inline-flex items-center gap-1.5 text-sm text-gray-500 dark:text-dark-400"
          >
            <Icon :name="tag.icon" size="sm" class="brand-text" />
            {{ tag.label }}
          </span>
        </div>
      </div>
    </main>

    <footer
      class="relative z-10 min-w-0 border-t border-gray-200/60 px-4 py-5 text-center text-sm text-gray-500 [overflow-wrap:anywhere] sm:px-6 dark:border-white/5 dark:text-dark-400"
    >
      &copy; {{ currentYear }} {{ siteName }}
    </footer>
  </div>

  <!-- Default Home Page -->
  <div
    v-else
    data-testid="default-home"
    :style="brandStyleVars"
    class="relative min-h-screen overflow-x-clip brand-wash text-gray-900 dark:text-white"
  >
    <!-- ==================== Header ==================== -->
    <header
      class="sticky top-0 z-40 border-b border-gray-200/60 bg-white/80 backdrop-blur-xl dark:border-white/5 dark:bg-dark-950/80"
    >
      <nav class="mx-auto flex max-w-7xl items-center justify-between gap-4 px-4 py-3 sm:px-6 lg:px-8">
        <!-- Brand -->
        <router-link to="/" class="flex min-w-0 items-center gap-2.5">
          <span
            class="flex h-9 w-9 shrink-0 items-center justify-center overflow-hidden rounded-xl brand-grad brand-shadow"
          >
            <img :src="siteLogo || '/logo.svg'" alt="Logo" class="h-full w-full object-contain" />
          </span>
          <span class="min-w-0 truncate text-base font-semibold tracking-tight">{{ siteName }}</span>
        </router-link>

        <!-- Section anchors -->
        <div class="hidden items-center gap-1 lg:flex">
          <a
            v-for="link in navLinks"
            :key="link.href"
            :href="link.href"
            class="rounded-lg px-3 py-2 text-sm font-medium text-gray-600 transition-colors hover:bg-gray-100 hover:text-gray-900 dark:text-dark-300 dark:hover:bg-white/5 dark:hover:text-white"
          >
            {{ link.label }}
          </a>
        </div>

        <!-- Actions -->
        <div class="flex shrink-0 items-center gap-1.5">
          <LocaleSwitcher />

          <a
            v-if="docUrl"
            :href="docUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="hidden h-9 w-9 items-center justify-center rounded-lg text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-900 sm:flex dark:text-dark-400 dark:hover:bg-white/5 dark:hover:text-white"
            :title="t('home.viewDocs')"
          >
            <Icon name="book" size="md" />
          </a>

          <router-link
            v-if="showModelPlazaEntry"
            to="/model-plaza"
            class="inline-flex h-9 items-center gap-1.5 rounded-lg px-2.5 text-sm font-medium text-gray-600 transition-colors hover:bg-gray-100 hover:text-gray-900 dark:text-dark-300 dark:hover:bg-white/5 dark:hover:text-white"
            :title="t('nav.modelPlaza')"
          >
            <Icon name="grid" size="md" />
            <span class="hidden md:inline">{{ t('nav.modelPlaza') }}</span>
          </router-link>

          <button
            class="flex h-9 w-9 items-center justify-center rounded-lg text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-900 dark:text-dark-400 dark:hover:bg-white/5 dark:hover:text-white"
            :title="isDark ? t('home.switchToLight') : t('home.switchToDark')"
            :aria-label="isDark ? t('home.switchToLight') : t('home.switchToDark')"
            @click="toggleTheme"
          >
            <Icon v-if="isDark" name="sun" size="md" />
            <Icon v-else name="moon" size="md" />
          </button>

          <!-- Signed in: avatar chip -->
          <router-link
            v-if="isAuthenticated"
            :to="dashboardPath"
            class="ml-1 inline-flex items-center gap-1.5 rounded-full bg-gray-900 py-1 pl-1 pr-2.5 transition-colors hover:bg-gray-800 dark:bg-white/10 dark:hover:bg-white/15"
          >
            <span
              class="flex h-6 w-6 items-center justify-center rounded-full brand-grad text-[10px] font-semibold text-white"
            >
              {{ userInitial }}
            </span>
            <span class="text-xs font-medium text-white">{{ t('home.dashboard') }}</span>
          </router-link>

          <!-- Signed out: primary CTA -->
          <router-link
            v-else
            to="/login"
            class="ml-1 rounded-full brand-grad brand-grad-hover brand-shadow px-4 py-2 text-sm font-medium text-white transition-all duration-200 active:scale-[0.98]"
          >
            {{ t('home.login') }}
          </router-link>
        </div>
      </nav>
    </header>

    <main class="relative">
      <!-- ==================== Hero ==================== -->
      <section class="relative overflow-hidden px-4 pb-20 pt-14 sm:px-6 lg:px-8 lg:pb-28 lg:pt-20">
        <!-- Ambient background -->
        <div class="pointer-events-none absolute inset-0" aria-hidden="true">
          <div
            class="absolute -right-24 -top-32 h-[26rem] w-[26rem] rounded-full brand-orb-from blur-3xl"
          ></div>
          <div
            class="absolute -left-32 top-24 h-[22rem] w-[22rem] rounded-full brand-orb-to blur-3xl"
          ></div>
          <div
            class="absolute bottom-0 left-1/2 h-[18rem] w-[34rem] -translate-x-1/2 rounded-full brand-orb-from-soft blur-3xl"
          ></div>
          <div class="hero-grid absolute inset-0"></div>
          <div
            class="absolute inset-x-0 top-0 h-64 bg-gradient-to-b from-white/70 to-transparent dark:from-dark-950/70"
          ></div>
        </div>

        <div
          class="relative mx-auto grid items-center gap-14"
          :class="
            showTerminal
              ? 'max-w-7xl lg:grid-cols-[1.05fr_1fr] lg:gap-16'
              : 'max-w-4xl lg:grid-cols-1'
          "
        >
          <!-- Left: copy -->
          <div class="text-center lg:text-left">
            <span
              class="inline-flex items-center gap-2 rounded-full brand-border brand-text border bg-white/70 px-3.5 py-1.5 text-xs font-medium shadow-sm backdrop-blur-sm dark:bg-white/5"
            >
              <span class="relative flex h-1.5 w-1.5">
                <span class="absolute inline-flex h-full w-full brand-bg animate-ping rounded-full opacity-75"></span>
                <span class="brand-bg relative inline-flex h-1.5 w-1.5 rounded-full"></span>
              </span>
              {{ t('home.heroBadge') }}
            </span>

            <h1
              class="mt-6 text-balance text-4xl font-bold leading-[1.1] tracking-tight sm:text-5xl lg:text-6xl"
            >
              {{ heroTitle }}
            </h1>

            <p
              class="mx-auto mt-6 max-w-xl text-pretty text-base leading-relaxed text-gray-600 sm:text-lg lg:mx-0 dark:text-dark-300"
            >
              {{ heroDescription }}
            </p>

            <!-- CTAs -->
            <div class="mt-9 flex flex-col items-center gap-3 sm:flex-row lg:justify-start">
              <router-link
                :to="isAuthenticated ? dashboardPath : '/login'"
                class="inline-flex w-full items-center justify-center gap-2 rounded-full brand-grad brand-grad-hover brand-shadow px-7 py-3 text-base font-medium text-white transition-all duration-200 active:scale-[0.98] sm:w-auto"
              >
                {{ isAuthenticated ? t('home.goToDashboard') : t('home.getStarted') }}
                <Icon name="arrowRight" size="md" :stroke-width="2" />
              </router-link>
              <a
                v-if="docUrl"
                :href="docUrl"
                target="_blank"
                rel="noopener noreferrer"
                class="btn btn-secondary btn-lg w-full rounded-full px-7 sm:w-auto"
              >
                <Icon name="book" size="md" />
                {{ t('home.viewDocs') }}
              </a>
            </div>

            <!-- Value tags -->
            <div class="mt-9 flex flex-wrap items-center justify-center gap-x-6 gap-y-3 lg:justify-start">
              <span
                v-for="tag in heroTags"
                :key="tag.label"
                class="inline-flex items-center gap-2 text-sm font-medium text-gray-600 dark:text-dark-300"
              >
                <Icon :name="tag.icon" size="sm" class="brand-text" />
                {{ tag.label }}
              </span>
            </div>
          </div>

          <!-- Right: request terminal -->
          <div v-if="showTerminal" class="flex justify-center lg:justify-end">
            <div class="terminal-container">
              <div class="terminal-window">
                <!-- Window header -->
                <div class="terminal-header">
                  <div class="terminal-buttons">
                    <span class="btn-close"></span>
                    <span class="btn-minimize"></span>
                    <span class="btn-maximize"></span>
                  </div>
                  <span class="terminal-title">gateway</span>
                </div>
                <!-- Terminal content -->
                <div class="terminal-body">
                  <div class="code-line line-1">
                    <span class="code-prompt">$</span>
                    <span class="code-cmd">curl</span>
                    <span class="code-flag">-X POST</span>
                    <span class="code-url">/v1/messages</span>
                  </div>
                  <div class="code-line line-2">
                    <span class="code-comment"># routing to the fastest upstream…</span>
                  </div>
                  <div class="code-line line-3">
                    <span class="code-success">200 OK</span>
                    <span class="code-response">{ "content": "Hello!" }</span>
                  </div>
                  <div class="code-line line-4">
                    <span class="code-dim">1.24s</span>
                    <span class="code-dim">·</span>
                    <span class="code-dim">$0.0031</span>
                    <span class="code-dim">·</span>
                    <span class="code-dim">claude-sonnet</span>
                  </div>
                  <div class="code-line line-5">
                    <span class="code-prompt">$</span>
                    <span class="cursor"></span>
                  </div>
                </div>
              </div>

              <!-- Floating caption -->
              <div class="terminal-badge">
                <Icon name="shield" size="sm" class="brand-text" />
                <span>{{ t('home.tags.stickySession') }}</span>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- ==================== Providers ==================== -->
      <section v-if="showProviders" id="providers" class="scroll-mt-24 px-4 py-14 sm:px-6 lg:px-8">
        <div class="mx-auto max-w-7xl">
          <div class="text-center">
            <h2 class="text-2xl font-bold tracking-tight sm:text-3xl">
              {{ t('home.providers.title') }}
            </h2>
            <p class="mt-3 text-sm text-gray-600 dark:text-dark-400">
              {{ t('home.providers.description') }}
            </p>
          </div>

          <div class="mt-10 grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-5">
            <div
              v-for="provider in providerList"
              :key="provider.model"
              class="group flex items-center gap-3 rounded-2xl border border-gray-200/70 bg-white/70 px-4 py-3.5 backdrop-blur-sm transition-all duration-300 brand-hover hover:-translate-y-0.5 dark:border-white/5 dark:bg-white/5"
            >
              <span
                class="flex h-9 w-9 shrink-0 items-center justify-center rounded-xl bg-white shadow-sm ring-1 ring-black/5 dark:ring-white/10"
              >
                <ModelIcon :model="provider.model" size="20px" />
              </span>
              <span class="min-w-0 truncate text-sm font-medium text-gray-700 dark:text-dark-200">
                {{ provider.label }}
              </span>
            </div>

            <!-- More -->
            <div
              class="flex items-center gap-3 rounded-2xl border border-dashed border-gray-300/70 bg-transparent px-4 py-3.5 dark:border-white/10"
            >
              <span
                class="flex h-9 w-9 shrink-0 items-center justify-center rounded-xl bg-gray-100 text-gray-500 dark:bg-white/5 dark:text-dark-400"
              >
                <Icon name="plus" size="sm" />
              </span>
              <span class="min-w-0 truncate text-sm font-medium text-gray-500 dark:text-dark-400">
                {{ t('home.providers.more') }}
              </span>
            </div>
          </div>
        </div>
      </section>

      <!-- ==================== Pain points ==================== -->
      <section v-if="showPainPoints" class="px-4 py-14 sm:px-6 lg:px-8">
        <div class="mx-auto max-w-7xl">
          <h2 class="text-center text-2xl font-bold tracking-tight sm:text-3xl">
            {{ t('home.painPoints.title') }}
          </h2>

          <div class="mt-10 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
            <div
              v-for="point in painPoints"
              :key="point.title"
              class="group rounded-2xl border border-gray-200/70 bg-white/70 p-6 backdrop-blur-sm transition-all duration-300 hover:border-rose-400/30 hover:shadow-lg hover:shadow-rose-500/5 dark:border-white/5 dark:bg-white/5"
            >
              <span
                class="flex h-11 w-11 items-center justify-center rounded-xl bg-rose-50 text-rose-500 transition-transform duration-300 group-hover:scale-110 dark:bg-rose-500/10 dark:text-rose-400"
              >
                <Icon :name="point.icon" size="lg" />
              </span>
              <h3 class="mt-4 text-base font-semibold">{{ point.title }}</h3>
              <p class="mt-2 text-sm leading-relaxed text-gray-600 dark:text-dark-400">
                {{ point.desc }}
              </p>
            </div>
          </div>
        </div>
      </section>

      <!-- ==================== Solutions / features ==================== -->
      <section id="features" class="scroll-mt-24 px-4 py-14 sm:px-6 lg:px-8">
        <div class="mx-auto max-w-7xl">
          <div class="text-center">
            <h2 class="text-2xl font-bold tracking-tight sm:text-3xl">
              {{ t('home.solutions.title') }}
            </h2>
            <p class="mt-3 text-sm text-gray-600 dark:text-dark-400">
              {{ t('home.solutions.subtitle') }}
            </p>
          </div>

          <div class="mt-12 grid gap-6 lg:grid-cols-3">
            <div
              v-for="(step, index) in solutionSteps"
              :key="step.title"
              class="group relative overflow-hidden rounded-3xl border border-gray-200/70 bg-white/70 p-7 backdrop-blur-sm transition-all duration-300 brand-hover hover:-translate-y-1 dark:border-white/5 dark:bg-white/5"
            >
              <!-- Step index watermark -->
              <span
                class="pointer-events-none absolute right-4 top-3 text-6xl font-bold leading-none text-gray-900/5 dark:text-white/5"
              >
                {{ String(index + 1).padStart(2, '0') }}
              </span>

              <span
                class="relative flex h-12 w-12 items-center justify-center rounded-2xl brand-grad brand-shadow text-white transition-transform duration-300 group-hover:scale-110"
              >
                <Icon :name="step.icon" size="lg" />
              </span>

              <h3 class="relative mt-5 text-lg font-semibold">{{ step.title }}</h3>
              <p class="relative mt-2.5 text-sm leading-relaxed text-gray-600 dark:text-dark-400">
                {{ step.desc }}
              </p>
            </div>
          </div>
        </div>
      </section>

      <!-- ==================== Comparison ==================== -->
      <section v-if="showComparison" id="comparison" class="scroll-mt-24 px-4 py-14 sm:px-6 lg:px-8">
        <div class="mx-auto max-w-4xl">
          <h2 class="text-center text-2xl font-bold tracking-tight sm:text-3xl">
            {{ t('home.comparison.title') }}
          </h2>

          <div
            class="mt-10 overflow-hidden rounded-3xl border border-gray-200/70 bg-white/70 shadow-sm backdrop-blur-sm dark:border-white/5 dark:bg-white/5"
          >
            <div class="overflow-x-auto">
              <div class="min-w-[620px]">
                <!-- Head -->
                <div class="grid grid-cols-[1.1fr_1.3fr_1.3fr] border-b border-gray-200/70 dark:border-white/5">
                  <div
                    class="px-5 py-4 text-xs font-semibold uppercase tracking-wider text-gray-500 dark:text-dark-400"
                  >
                    {{ t('home.comparison.headers.feature') }}
                  </div>
                  <div
                    class="px-5 py-4 text-xs font-semibold uppercase tracking-wider text-gray-500 dark:text-dark-400"
                  >
                    {{ t('home.comparison.headers.official') }}
                  </div>
                  <div
                    class="flex items-center gap-2 brand-soft-strong brand-text px-5 py-4 text-xs font-semibold uppercase tracking-wider"
                  >
                    {{ t('home.comparison.headers.us') }}
                    <span
                      class="brand-bg rounded-full px-2 py-0.5 text-[10px] font-semibold normal-case text-white"
                    >
                      {{ t('home.comparison.badge') }}
                    </span>
                  </div>
                </div>

                <!-- Rows -->
                <div
                  v-for="(row, index) in comparisonRows"
                  :key="row.feature"
                  class="grid grid-cols-[1.1fr_1.3fr_1.3fr]"
                  :class="index !== comparisonRows.length - 1
                    ? 'border-b border-gray-200/70 dark:border-white/5'
                    : ''"
                >
                  <div class="px-5 py-4 text-sm font-medium text-gray-900 dark:text-white">
                    {{ row.feature }}
                  </div>
                  <div class="px-5 py-4 text-sm text-gray-500 dark:text-dark-400">
                    {{ row.official }}
                  </div>
                  <div class="flex items-start gap-2 brand-soft px-5 py-4 text-sm font-medium text-gray-900 dark:text-white">
                    <Icon name="check" size="sm" class="brand-text mt-0.5 shrink-0" :stroke-width="2.5" />
                    {{ row.us }}
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- ==================== CTA ==================== -->
      <section class="px-4 py-14 sm:px-6 lg:px-8">
        <div class="mx-auto max-w-5xl">
          <div
            class="relative overflow-hidden rounded-3xl brand-grad brand-shadow-lg px-6 py-14 text-center sm:px-12"
          >
            <div class="pointer-events-none absolute inset-0" aria-hidden="true">
              <div class="absolute -right-16 -top-16 h-56 w-56 rounded-full bg-white/15 blur-3xl"></div>
              <div class="absolute -bottom-20 -left-10 h-56 w-56 brand-orb-to rounded-full blur-3xl"></div>
              <div class="cta-grid absolute inset-0 opacity-40"></div>
            </div>

            <div class="relative">
              <h2 class="text-2xl font-bold tracking-tight text-white sm:text-3xl">
                {{ t('home.cta.title') }}
              </h2>
              <p class="mx-auto mt-4 max-w-xl text-pretty text-sm leading-relaxed text-white/85 sm:text-base">
                {{ t('home.cta.description') }}
              </p>
              <router-link
                :to="isAuthenticated ? dashboardPath : '/register'"
                class="mt-8 inline-flex items-center gap-2 rounded-full brand-cta-btn bg-white px-7 py-3 text-sm font-semibold shadow-lg shadow-black/10 transition-transform duration-200 hover:-translate-y-0.5"
              >
                {{ isAuthenticated ? t('home.goToDashboard') : t('home.cta.button') }}
                <Icon name="arrowRight" size="md" :stroke-width="2" />
              </router-link>
            </div>
          </div>
        </div>
      </section>
    </main>

    <!-- ==================== Footer ==================== -->
    <footer class="border-t border-gray-200/70 px-4 py-10 sm:px-6 lg:px-8 dark:border-white/5">
      <div class="mx-auto flex max-w-7xl flex-col items-center justify-between gap-5 sm:flex-row">
        <div class="flex items-center gap-2.5">
          <span
            class="flex h-8 w-8 items-center justify-center overflow-hidden rounded-lg brand-grad"
          >
            <img :src="siteLogo || '/logo.svg'" alt="Logo" class="h-full w-full object-contain" />
          </span>
          <span class="text-sm font-medium">{{ siteName }}</span>
        </div>

        <div class="flex flex-wrap items-center justify-center gap-x-6 gap-y-2">
          <a
            v-if="docUrl"
            :href="docUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="text-sm text-gray-500 transition-colors hover:text-gray-900 dark:text-dark-400 dark:hover:text-white"
          >
            {{ t('home.docs') }}
          </a>
          <router-link
            v-if="showModelPlazaEntry"
            to="/model-plaza"
            class="text-sm text-gray-500 transition-colors hover:text-gray-900 dark:text-dark-400 dark:hover:text-white"
          >
            {{ t('nav.modelPlaza') }}
          </router-link>
          <a
            :href="githubUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="text-sm text-gray-500 transition-colors hover:text-gray-900 dark:text-dark-400 dark:hover:text-white"
          >
            GitHub
          </a>
        </div>

        <p class="text-sm text-gray-500 dark:text-dark-400">
          &copy; {{ currentYear }} {{ siteName }}. {{ t('home.footer.allRightsReserved') }}
        </p>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore, useAppStore } from '@/stores'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import ModelIcon from '@/components/common/ModelIcon.vue'
import Icon from '@/components/icons/Icon.vue'
import { sanitizeUrl } from '@/utils/url'
import { FeatureFlags, isFeatureFlagEnabled } from '@/utils/featureFlags'
import { homeStyleCssVars, normalizeHomeStyle } from '@/constants/homeStyle'

// Subset of Icon names used on this page; keeps v-for data type-safe.
type IconName =
  | 'creditCard'
  | 'users'
  | 'exclamationTriangle'
  | 'chartBar'
  | 'key'
  | 'swap'
  | 'shield'
  | 'chart'
  | 'check'
  | 'arrowRight'
  | 'book'
  | 'grid'
  | 'sun'
  | 'moon'
  | 'plus'

const { t } = useI18n()

const authStore = useAuthStore()
const appStore = useAppStore()

// Site settings - directly from appStore (already initialized from injected config)
const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || 'Sub2API')
const siteLogo = computed(() => sanitizeUrl(appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '', { allowRelative: true, allowDataUrl: true }))
const siteSubtitle = computed(() => appStore.cachedPublicSettings?.site_subtitle || 'AI API Gateway Platform')
const docUrl = computed(() => sanitizeUrl(appStore.cachedPublicSettings?.doc_url || appStore.docUrl || ''))
const homeContent = computed(() => appStore.cachedPublicSettings?.home_content || '')
const hasHomeContent = computed(() => homeContent.value.trim().length > 0)
const compactHomeEnabled = computed(() => appStore.cachedPublicSettings?.compact_home_enabled === true)
const modelPlazaEnabled = computed(() => isFeatureFlagEnabled(FeatureFlags.modelPlaza))

// Admin-configurable homepage appearance. normalizeHomeStyle() fills in the
// shipped defaults so an older or partial payload cannot blank the page.
const homeStyle = computed(() => normalizeHomeStyle(appStore.cachedPublicSettings?.home_style))
const brandStyleVars = computed(() => homeStyleCssVars(homeStyle.value))
const heroTitle = computed(() => homeStyle.value.hero_title || t('home.heroSubtitle'))
const heroDescription = computed(() => homeStyle.value.hero_desc || t('home.heroDescription'))
const showProviders = computed(() => homeStyle.value.show_providers)
const showPainPoints = computed(() => homeStyle.value.show_pain_points)
const showComparison = computed(() => homeStyle.value.show_comparison)
const showTerminal = computed(() => homeStyle.value.show_terminal)

// Check if homeContent is a URL (for iframe display)
const isHomeContentUrl = computed(() => {
  const content = homeContent.value.trim()
  return content.startsWith('http://') || content.startsWith('https://')
})

// Theme
const isDark = ref(document.documentElement.classList.contains('dark'))

// GitHub URL
const githubUrl = 'https://github.com/Wei-Shaw/sub2api'

// Auth state
const isAuthenticated = computed(() => authStore.isAuthenticated)
const modelPlazaRequiresAuth = computed(
  () => appStore.cachedPublicSettings?.model_plaza_require_auth === true,
)
const showModelPlazaEntry = computed(
  () => modelPlazaEnabled.value && (isAuthenticated.value || !modelPlazaRequiresAuth.value),
)
const isAdmin = computed(() => authStore.isAdmin)
const dashboardPath = computed(() => isAdmin.value ? '/admin/dashboard' : '/dashboard')
const userInitial = computed(() => {
  const user = authStore.user
  if (!user || !user.email) return ''
  return user.email.charAt(0).toUpperCase()
})

// Current year for footer
const currentYear = computed(() => new Date().getFullYear())

// Page content ------------------------------------------------------------
// Anchors for sections that are currently visible; a hidden section must not
// leave a dead link in the header.
const navLinks = computed(() => {
  const links: { href: string; label: string }[] = [{ href: '#features', label: t('home.nav.features') }]
  if (showProviders.value) {
    links.push({ href: '#providers', label: t('home.nav.providers') })
  }
  if (showComparison.value) {
    links.push({ href: '#comparison', label: t('home.nav.comparison') })
  }
  return links
})

const heroTags = computed<{ icon: IconName; label: string }[]>(() => [
  { icon: 'swap', label: t('home.tags.subscriptionToApi') },
  { icon: 'shield', label: t('home.tags.stickySession') },
  { icon: 'chart', label: t('home.tags.realtimeBilling') },
])

// Brand marks come from ModelIcon's bundled @lobehub icon set.
const providerList = computed(() => [
  { model: 'claude', label: t('home.providers.claude') },
  { model: 'openai', label: t('home.providers.openai') },
  { model: 'gemini', label: t('home.providers.gemini') },
  { model: 'deepseek', label: t('home.providers.deepseek') },
  { model: 'qwen', label: t('home.providers.qwen') },
  { model: 'xai', label: t('home.providers.grok') },
  { model: 'moonshot', label: t('home.providers.kimi') },
  { model: 'zhipu', label: t('home.providers.zhipu') },
  { model: 'antigravity', label: t('home.providers.antigravity') },
])

const painPoints = computed<{ icon: IconName; title: string; desc: string }[]>(() => [
  {
    icon: 'creditCard',
    title: t('home.painPoints.items.expensive.title'),
    desc: t('home.painPoints.items.expensive.desc'),
  },
  {
    icon: 'users',
    title: t('home.painPoints.items.complex.title'),
    desc: t('home.painPoints.items.complex.desc'),
  },
  {
    icon: 'exclamationTriangle',
    title: t('home.painPoints.items.unstable.title'),
    desc: t('home.painPoints.items.unstable.desc'),
  },
  {
    icon: 'chartBar',
    title: t('home.painPoints.items.noControl.title'),
    desc: t('home.painPoints.items.noControl.desc'),
  },
])

const solutionSteps = computed<{ icon: IconName; title: string; desc: string }[]>(() => [
  {
    icon: 'key',
    title: t('home.features.unifiedGateway'),
    desc: t('home.features.unifiedGatewayDesc'),
  },
  {
    icon: 'swap',
    title: t('home.features.multiAccount'),
    desc: t('home.features.multiAccountDesc'),
  },
  {
    icon: 'creditCard',
    title: t('home.features.balanceQuota'),
    desc: t('home.features.balanceQuotaDesc'),
  },
])

const comparisonRows = computed(() => [
  {
    feature: t('home.comparison.items.pricing.feature'),
    official: t('home.comparison.items.pricing.official'),
    us: t('home.comparison.items.pricing.us'),
  },
  {
    feature: t('home.comparison.items.models.feature'),
    official: t('home.comparison.items.models.official'),
    us: t('home.comparison.items.models.us'),
  },
  {
    feature: t('home.comparison.items.management.feature'),
    official: t('home.comparison.items.management.official'),
    us: t('home.comparison.items.management.us'),
  },
  {
    feature: t('home.comparison.items.stability.feature'),
    official: t('home.comparison.items.stability.official'),
    us: t('home.comparison.items.stability.us'),
  },
  {
    feature: t('home.comparison.items.control.feature'),
    official: t('home.comparison.items.control.official'),
    us: t('home.comparison.items.control.us'),
  },
])

// Toggle theme
function toggleTheme() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

// Initialize theme
function initTheme() {
  const savedTheme = localStorage.getItem('theme')
  if (
    savedTheme === 'dark' ||
    (!savedTheme && window.matchMedia('(prefers-color-scheme: dark)').matches)
  ) {
    isDark.value = true
    document.documentElement.classList.add('dark')
  }
}

onMounted(() => {
  initTheme()

  // Check auth state
  authStore.checkAuth()

  // Ensure public settings are loaded (will use cache if already loaded from injected config)
  if (!appStore.publicSettingsLoaded) {
    appStore.fetchPublicSettings()
  }
})
</script>

<style scoped>
/* ============ Brand accent (admin-configurable via home_style) ============
   Every brand-colored surface reads these custom properties, so changing the
   accent in the admin settings re-tints the whole page without a rebuild.
   The *-rgb variables hold a space-separated triplet for `rgb(var(--x) / a)`. */
.brand-grad {
  background-image: linear-gradient(135deg, rgb(var(--brand-from-rgb)), rgb(var(--brand-to-rgb)));
}

.brand-grad-hover:hover {
  filter: brightness(1.08);
}

.brand-shadow {
  box-shadow: 0 10px 30px rgb(var(--brand-from-rgb) / 0.25);
}

.brand-shadow-lg {
  box-shadow: 0 25px 50px -12px rgb(var(--brand-from-rgb) / 0.35);
}

.brand-bg {
  background-color: rgb(var(--brand-from-rgb));
}

.brand-text {
  color: rgb(var(--brand-from-rgb));
}

.brand-border {
  border-color: rgb(var(--brand-from-rgb) / 0.25);
}

.brand-soft {
  background-color: rgb(var(--brand-from-rgb) / 0.05);
}

.brand-soft-strong {
  background-color: rgb(var(--brand-from-rgb) / 0.1);
}

.brand-orb-from {
  background-color: rgb(var(--brand-from-rgb) / 0.22);
}

.brand-orb-from-soft {
  background-color: rgb(var(--brand-from-rgb) / 0.14);
}

.brand-orb-to {
  background-color: rgb(var(--brand-to-rgb) / 0.18);
}

.brand-hover:hover {
  border-color: rgb(var(--brand-from-rgb) / 0.32);
  box-shadow: 0 10px 30px rgb(var(--brand-from-rgb) / 0.1);
}

.brand-cta-btn {
  color: rgb(var(--brand-from-rgb));
}

.brand-cta-btn:hover {
  background-color: rgb(var(--brand-from-rgb) / 0.08);
}

.brand-wash {
  background-image: linear-gradient(
    to bottom,
    #ffffff,
    rgb(var(--brand-from-rgb) / 0.06),
    #ffffff
  );
}

.dark .brand-wash {
  background-image: linear-gradient(to bottom, #020617, #0f172a, #020617);
}

/* ============ Background patterns ============ */
.hero-grid {
  background-image:
    linear-gradient(to right, rgb(var(--brand-from-rgb) / 0.08) 1px, transparent 1px),
    linear-gradient(to bottom, rgb(var(--brand-from-rgb) / 0.08) 1px, transparent 1px);
  background-size: 64px 64px;
  mask-image: radial-gradient(ellipse 90% 65% at 50% 0%, #000 30%, transparent 100%);
  -webkit-mask-image: radial-gradient(ellipse 90% 65% at 50% 0%, #000 30%, transparent 100%);
}

.dark .hero-grid {
  background-image:
    linear-gradient(to right, rgb(var(--brand-from-rgb) / 0.1) 1px, transparent 1px),
    linear-gradient(to bottom, rgb(var(--brand-from-rgb) / 0.1) 1px, transparent 1px);
}

.cta-grid {
  background-image:
    linear-gradient(to right, rgba(255, 255, 255, 0.14) 1px, transparent 1px),
    linear-gradient(to bottom, rgba(255, 255, 255, 0.14) 1px, transparent 1px);
  background-size: 44px 44px;
  mask-image: radial-gradient(ellipse 70% 80% at 50% 50%, #000 20%, transparent 90%);
  -webkit-mask-image: radial-gradient(ellipse 70% 80% at 50% 50%, #000 20%, transparent 90%);
}

/* ============ Terminal ============ */
.terminal-container {
  position: relative;
  display: inline-block;
}

.terminal-window {
  width: 100%;
  max-width: 460px;
  background: linear-gradient(145deg, #1e293b 0%, #0f172a 100%);
  border-radius: 16px;
  box-shadow:
    0 25px 50px -12px rgba(0, 0, 0, 0.4),
    0 0 0 1px rgba(255, 255, 255, 0.1),
    inset 0 1px 0 rgba(255, 255, 255, 0.1);
  overflow: hidden;
  transform: perspective(1000px) rotateX(2deg) rotateY(-2deg);
  transition: transform 0.3s ease;
}

.terminal-window:hover {
  transform: perspective(1000px) rotateX(0deg) rotateY(0deg) translateY(-4px);
}

.terminal-header {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  background: rgba(30, 41, 59, 0.8);
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.terminal-buttons {
  display: flex;
  gap: 8px;
}

.terminal-buttons span {
  width: 12px;
  height: 12px;
  border-radius: 50%;
}

.btn-close {
  background: #ef4444;
}
.btn-minimize {
  background: #eab308;
}
.btn-maximize {
  background: #22c55e;
}

.terminal-title {
  flex: 1;
  text-align: center;
  font-size: 12px;
  font-family: ui-monospace, monospace;
  color: #64748b;
  margin-right: 52px;
}

.terminal-body {
  padding: 20px 24px;
  font-family: ui-monospace, 'Fira Code', monospace;
  font-size: 13.5px;
  line-height: 2;
}

.code-line {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  opacity: 0;
  animation: line-appear 0.5s ease forwards;
}

.line-1 {
  animation-delay: 0.3s;
}
.line-2 {
  animation-delay: 1s;
}
.line-3 {
  animation-delay: 1.8s;
}
.line-4 {
  animation-delay: 2.4s;
}
.line-5 {
  animation-delay: 2.9s;
}

@keyframes line-appear {
  from {
    opacity: 0;
    transform: translateY(5px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.code-prompt {
  color: #22c55e;
  font-weight: bold;
}
.code-cmd {
  color: #38bdf8;
}
.code-flag {
  color: #a78bfa;
}
.code-url {
  color: #14b8a6;
}
.code-comment {
  color: #64748b;
  font-style: italic;
}
.code-success {
  color: #22c55e;
  background: rgba(34, 197, 94, 0.15);
  padding: 2px 8px;
  border-radius: 4px;
  font-weight: 600;
}
.code-response {
  color: #fbbf24;
}
.code-dim {
  color: #64748b;
  font-size: 12.5px;
}

/* Blinking cursor */
.cursor {
  display: inline-block;
  width: 8px;
  height: 16px;
  background: #22c55e;
  animation: blink 1s step-end infinite;
}

@keyframes blink {
  0%,
  50% {
    opacity: 1;
  }
  51%,
  100% {
    opacity: 0;
  }
}

/* Floating caption on the terminal */
.terminal-badge {
  position: absolute;
  bottom: -18px;
  left: -14px;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 8px 14px;
  border-radius: 9999px;
  font-size: 12px;
  font-weight: 500;
  color: #334155;
  background: rgba(255, 255, 255, 0.92);
  border: 1px solid rgba(15, 23, 42, 0.06);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
  backdrop-filter: blur(8px);
  animation: badge-float 4s ease-in-out infinite;
}

.dark .terminal-badge {
  color: #e2e8f0;
  background: rgba(30, 41, 59, 0.9);
  border-color: rgba(255, 255, 255, 0.08);
}

@keyframes badge-float {
  0%,
  100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-6px);
  }
}

/* Dark mode terminal glow */
.dark .terminal-window {
  box-shadow:
    0 25px 50px -12px rgba(0, 0, 0, 0.6),
    0 0 0 1px rgb(var(--brand-from-rgb) / 0.25),
    0 0 40px rgb(var(--brand-to-rgb) / 0.12),
    inset 0 1px 0 rgba(255, 255, 255, 0.1);
}

@media (max-width: 640px) {
  .terminal-badge {
    left: 4px;
    bottom: -16px;
  }
}

@media (prefers-reduced-motion: reduce) {
  .code-line,
  .cursor,
  .terminal-badge {
    animation: none;
    opacity: 1;
  }
  .terminal-window,
  .terminal-window:hover {
    transition: none;
    transform: none;
  }
}
</style>
