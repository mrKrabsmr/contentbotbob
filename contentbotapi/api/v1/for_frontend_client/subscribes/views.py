import http.client

from rest_framework.permissions import IsAuthenticated
from rest_framework.response import Response
from rest_framework.views import APIView

from apps.subscribes.serializers import SubscribeSerializer
from apps.subscribes.services import SubscribeService as S_s, UserSubscribeService as US_s
from apps.users.permissions import IsConfirmed


class SubscribeAPIView(APIView):
    permission_classes = [IsAuthenticated, IsConfirmed]

    def get(self, request, *args, **kwargs):
        subscribes = S_s.get_all()
        serializer = SubscribeSerializer(subscribes, many=True)
        return Response({"result": serializer.data}, http.HTTPStatus.OK)

    def post(self, request, *args, **kwargs):
        user = request.user
        sub_id = request.data.get("subscribe_id")
        sub = S_s.get_subscribe(sub_id)
        if not sub:
            return Response({"message": "подписка не найдена"}, http.HTTPStatus.NOT_FOUND)

        if US_s.check_user_already_have_subscibe(user):
            return Response({"message": "у вас уже есть подписка"}, http.HTTPStatus.BAD_REQUEST)

        US_s.activate_user_subscribe(user, sub)

        return Response({"message": "подписка уже активирована"}, http.HTTPStatus.CREATED)
