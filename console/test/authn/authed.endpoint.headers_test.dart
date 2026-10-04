import 'package:fixnum/fixnum.dart' as fixnum;
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/meta.dart' as meta;
import 'package:retrovibed/testing/widget_tester_extensions.dart';

void main() {
  group('AuthedEndpoint.headers', () {
    testWidgets('uses the scoped token for uris on the scoped daemon', (tester) async {
      Map<String, String>? headers;

      await tester.pumpApp(
        authn.AuthedEndpoint(
          initial: meta.Daemon(hostname: 'remote:1'),
          current: ({String? host}) async => meta.AuthzResponse(
            bearer: 'bearer abc',
            token: meta.Token(exp: fixnum.Int64(DateTime.now().millisecondsSinceEpoch ~/ 1000 + 3600)),
          ),
          Builder(
            builder: (context) {
              headers = authn.AuthedEndpoint.headers(context, 'https://remote:1/m/x');
              return const SizedBox();
            },
          ),
        ),
      );
      await tester.pumpN(5);

      expect(headers, {"Authorization": "bearer abc"});
      expect(tester.takeException(), isNull);
    });

    testWidgets('prefixes an unprefixed scoped token', (tester) async {
      Map<String, String>? headers;

      await tester.pumpApp(
        authn.AuthedEndpoint(
          initial: meta.Daemon(hostname: 'remote:1'),
          current: ({String? host}) async => meta.AuthzResponse(
            bearer: 'abc',
            token: meta.Token(exp: fixnum.Int64(DateTime.now().millisecondsSinceEpoch ~/ 1000 + 3600)),
          ),
          Builder(
            builder: (context) {
              headers = authn.AuthedEndpoint.headers(context, 'https://remote:1/m/x');
              return const SizedBox();
            },
          ),
        ),
      );
      await tester.pumpN(5);

      expect(headers, {"Authorization": "bearer abc"});
      expect(tester.takeException(), isNull);
    });

    testWidgets('falls back to local headers for uris on another host', (tester) async {
      Map<String, String>? headers = const {};

      await tester.pumpApp(
        authn.AuthedEndpoint(
          initial: meta.Daemon(hostname: 'remote:1'),
          current: ({String? host}) async => meta.AuthzResponse(
            bearer: 'bearer abc',
            token: meta.Token(exp: fixnum.Int64(DateTime.now().millisecondsSinceEpoch ~/ 1000 + 3600)),
          ),
          Builder(
            builder: (context) {
              headers = authn.AuthedEndpoint.headers(context, 'https://other:2/m/x');
              return const SizedBox();
            },
          ),
        ),
      );
      await tester.pumpN(5);

      expect(headers, httpx.localheaders('https://other:2/m/x'));
      expect(tester.takeException(), isNull);
    });

    testWidgets('falls back to local headers without an AuthedEndpoint ancestor', (tester) async {
      Map<String, String>? headers = const {};

      await tester.pumpApp(
        Builder(
          builder: (context) {
            headers = authn.AuthedEndpoint.headers(context, 'https://remote:1/m/x');
            return const SizedBox();
          },
        ),
      );
      await tester.pumpN(5);

      expect(headers, httpx.localheaders('https://remote:1/m/x'));
      expect(tester.takeException(), isNull);
    });
  });
}
