<template>
    <div>
        <button class="btn" @click="openModal">open modal</button>

        <dialog ref="modal" class="modal">
            <div class="modal-box">
                <h3 class="text-lg font-bold">{{ title }}</h3>
                <p class="py-4">{{ content }}</p>
                <div class="modal-action">
                    <form method="dialog">
                        <button class="btn" @click="closeModal">Close</button>
                    </form>
                </div>
            </div>
        </dialog>
    </div>
</template>

<script>
export default {
    props: {
        title: {
            type: String,
            default: "Hello!"
        },
        content: {
            type: String,
            default: "Press ESC key or click the button below to close"
        }
    },
    methods: {
        openModal() {
            this.$refs.modal.showModal();
        },
        closeModal() {
            this.$refs.modal.close();
        }
    }
};
</script>