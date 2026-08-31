import 'package:flutter/material.dart';
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/billing/api.dart' as billing;
import 'package:retrovibed/community/api.dart' as api;
import 'package:retrovibed/community/community.pb.dart';

Future<void> handleSubscribeAction(BuildContext context, Community community, String attribution) {
  final auth = [authn.request(authn.AuthzCache.meta(context))];
  billing.consumeAttribution(attribution, options: [authn.Authenticated.bearer(context)]).ignore();
  return httpx.withRetry(() => api.communities.subscribe(community.id, options: auth));
}
