import { PUBLIC_PB_ADDR } from '$env/static/public';
import PocketBase from 'pocketbase';

export const pb = new PocketBase(PUBLIC_PB_ADDR);
export const PB_COLLECTION_USERS = 'users';
export const PB_COLLECTION_UPSCALE_RUNNERS = 'upscale_runners';
export const PB_COLLECTION_IMAGE_QUEUES = 'image_queues';
export const PB_COLLECTION_GENERATED_IMAGES = 'generated_images';
export const PB_COLLECTION_SETTINGS = 'settings';
export const PB_COLLECTION_UPACALE_QUEUES = 'upscale_queues';
