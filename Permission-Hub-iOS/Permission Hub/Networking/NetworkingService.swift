//
//  NetworkingService.swift
//  Permission Hub
//
//  Created by Corné on 12/04/2019.
//  Copyright © 2019 Planet. All rights reserved.
//

import Foundation

enum PHNetworkingError: Error {
    case invalidUrl
}

extension PHNetworkingError: LocalizedError {

    var errorDescription: String? {
        switch self {
        case .invalidUrl:
            return "unable to construct URL"
        }
    }
}

enum PHNetworkingResult<T> {
    case success(T)
    case failure(Error)
}

final class NetworkingService {

    // MARK: - Private properties

    private let baseUrl = "http://51.15.52.136"

    // MARK: - Endpoints

    // endpoint to poll for notifications
    func getNotifications(completion: @escaping  (_ result: PHNetworkingResult<TransactionNotification>) -> Void) throws {

        // construct URL
        guard let url = URL(string: baseUrl + "/notification-get") else {
            throw PHNetworkingError.invalidUrl
        }

        let task = URLSession.shared.dataTask(with: url) { data, response, error in

            guard let response = response as? HTTPURLResponse,
                response.statusCode == 500 else {
                    return
            }

            if let error = error {
                DispatchQueue.main.async {
                    completion(.failure(error))
                }
            }

            if let data = data {

                let decoder = JSONDecoder()
                decoder.dateDecodingStrategy = .iso8601

                do {
                    let notification = try decoder.decode(TransactionNotification.self, from: data)
                    DispatchQueue.main.async {
                        completion(.success(notification))
                    }
                } catch {
                    DispatchQueue.main.async {
                        completion(.failure(error))
                    }
                }
            }
        }

        task.resume()
    }
}
