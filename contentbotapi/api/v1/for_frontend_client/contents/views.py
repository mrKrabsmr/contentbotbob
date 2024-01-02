import http

from rest_framework.response import Response
from rest_framework.views import APIView

from apps.contents.serializers import ContentSerializer
from apps.contents.services import ContentService as C_s

from config.settings import CLIENT_KEY


class ContentAPIView(APIView):

    def get(self, request, outer_id, *args, **kwargs):
        client_key = request.query_params.get("key")
        if client_key != CLIENT_KEY:
            return Response({"message": "incorrect key"}, http.HTTPStatus.FORBIDDEN)

        content, active = C_s.get_suitable_one(outer_id)
        if not content:
            return Response({"result": None, "active": active}, http.HTTPStatus.NOT_FOUND)

        serializer = ContentSerializer(content)
        return Response({"result": serializer.data, "active": active}, http.HTTPStatus.OK)
