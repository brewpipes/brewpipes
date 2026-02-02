<template>
  <v-container class="production-page" fluid>
    <v-row>
      <v-col cols="12">
        <v-card class="section-card">
          <v-card-title>
            <v-icon class="mr-2" icon="mdi-factory" />
            Production
          </v-card-title>
          <v-card-text>
            <v-row>
              <v-col cols="12" md="4">
                <v-card to="/production/recipes" variant="tonal">
                  <v-card-title>
                    <v-icon class="mr-2" icon="mdi-book-open-page-variant" />
                    Recipes
                  </v-card-title>
                  <v-card-text>
                    Manage beer recipes and styles for your brewery.
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts" setup>
</script>

<style scoped>
.production-page {
  position: relative;
}

.section-card {
  background: rgba(var(--v-theme-surface), 0.92);
  border: 1px solid rgba(var(--v-theme-on-surface), 0.12);
  box-shadow: 0 12px 26px rgba(0, 0, 0, 0.2);
}
</style>
