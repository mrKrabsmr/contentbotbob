import logging

from django.db.transaction import atomic

from apps.contentmaker.models import Channel, ChannelContent, ChannelContentMedia


def filter_channels(data):
    empty = False if data.get("img") else True

    channels = Channel.objects.filter(
        settings__allowed_content_sources__value=data.get("source"),
        settings__min_rating__lte=data.get("rating")
    )

    if empty:
        channels = channels.filter(settings__empty_file_allowed=True)

    return channels


def evaluate_content(data):
    rating = 0

    if data.get("img"):
        rating += 3

    len_text = len(data.get("text"))
    if len_text < 50 or len_text > 1500:
        rating += 1
    elif len_text < 100 or len_text > 1000:
        rating += 2
    elif len_text < 200 or len_text > 900:
        rating += 3
    elif len_text < 300 or len_text > 800:
        rating += 4
    elif len_text < 350 or len_text > 750:
        rating += 5
    elif len_text < 400 or len_text > 700:
        rating += 6
    else:
        rating += 7

    return rating


@atomic
def saving_data(data):
    data["rating"] = evaluate_content(data)

    contents = []
    content_images = []

    for channel in filter_channels(data):
        content = ChannelContent(text=data["text"], rating=data["rating"], channel=channel)
        content_image = ChannelContentMedia(file_url=data["img"], channel_content=content)

        contents.append(content)
        content_images.append(content_image)

    ChannelContent.objects.bulk_create(contents)
    ChannelContentMedia.objects.bulk_create(content_images)
