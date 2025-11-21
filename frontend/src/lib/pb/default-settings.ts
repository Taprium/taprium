import { pb, PB_COLLECTION_SETTINGS } from './backend-pb';

export async function GetDefaultSettings() {
	return await pb.collection(PB_COLLECTION_SETTINGS).getFirstListItem('');
}
