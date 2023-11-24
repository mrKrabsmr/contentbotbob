import http

from rest_framework.response import Response
from rest_framework.views import APIView

from apps.contents.serializers import NewsContentSerializer
from apps.contents.services import ContentDistributionService as CD_s


class NewsContentAPIView(APIView):
    def post(self, request, *args, **kwargs):
        serializer = NewsContentSerializer(data=request.data, many=True)
        if serializer.is_valid():
            for data in serializer.data:
                CD_s.saving_data(data)
            return Response({"message": "success"}, http.HTTPStatus.CREATED)
        return Response(serializer.errors, http.HTTPStatus.BAD_REQUEST)
