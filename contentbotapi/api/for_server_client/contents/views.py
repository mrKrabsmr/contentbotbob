import logging

from rest_framework.response import Response
from rest_framework.views import APIView

from apps.contentmaker.serializers import NewsContentSerializer
from apps.contentmaker.utils import saving_data


class NewsContentAPIView(APIView):
    def post(self, request, *args, **kwargs):
        serializer = NewsContentSerializer(data=request.data, many=True)
        if serializer.is_valid():
            for data in serializer.data:
                saving_data(data)
            return Response({"message": "success"}, 201)
        return Response(serializer.errors, 400)
