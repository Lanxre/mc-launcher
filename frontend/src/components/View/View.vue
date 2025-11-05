<template>
    <div class="main-content">
        <slot></slot>
    </div>
</template>

<style lang="css" scoped>
.main-content {
  display: flex;
  flex-direction: column;
  
  width: 100%;
  
  padding: 20px;
  
  gap: 15px;
}
</style>