import 'dart:async';
import 'dart:convert';
import 'package:retrovibed/audio/meta.audio.pb.dart';
import 'package:retrovibed/httpx.dart' as httpx;

export 'package:retrovibed/audio/meta.audio.pb.dart';

abstract class sinks {
  static Future<Stream<AudioSink>> listen({
    List<httpx.Option> options = const [],
  }) async {
    return httpx.websocket(Uri.https(httpx.host(), "/audio/sinks/", null), options: options).then((socket) {
      return socket.transform(
        StreamTransformer.fromHandlers(
          handleData: (data, sink) {
            if (data is List<int>) {
              final resp = httpx.fromProto3JsonSafe(AudioSinkSearchResponse.create(), jsonDecode(utf8.decode(data)));
              resp.items.forEach(sink.add);
            } else {
              sink.addError('deserialization failed data: $data');
            }
          },
        ),
      );
    });
  }

  static Future<AudioSinkCurrentResponse> current({
    List<httpx.Option> options = const [],
  }) async {
    return httpx.get(Uri.https(httpx.host(), "/audio/sinks/", {}), options: options).then((v) {
      return Future.value(
        httpx.fromProto3JsonSafe(AudioSinkCurrentResponse.create(), jsonDecode(v.body)),
      );
    });
  }

  static Future<AudioSinkTouchResponse> activate(
    String id, {
    List<httpx.Option> options = const [],
  }) async {
    return httpx
        .post(
          Uri.https(httpx.host(), "/audio/sinks/", {}),
          options: options,
          body: jsonEncode(AudioSinkTouchRequest(id: id).toProto3Json()),
        )
        .then((v) {
          return Future.value(
            httpx.fromProto3JsonSafe(AudioSinkTouchResponse.create(), jsonDecode(v.body)),
          );
        });
  }
}
