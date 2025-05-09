<template>
  <header class="navbar bg-base-100 shadow-xl fixed z-50">
    <nav v-if="$slots.start" class="navbar-start ml-2">
      <slot name="start"></slot>
    </nav>
    <nav v-if="$slots.center" class="navbar-center">
      <slot name="center"></slot>
    </nav>
    <nav v-if="$slots.end" class="navbar-end mr-2">
      <slot name="end"></slot>
    </nav>
  </header>
</template>
