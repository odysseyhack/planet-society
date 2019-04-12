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

final class NetworkingService {

    // MARK: - Private properties

    private let baseUrl = "http://51.15.52.136"

    // MARK: - Endpoints

    // endpoint to poll for notifications
    func getNotifications() throws {

        // construct URL
        guard let url = URL(string: baseUrl + "/notification-get") else {
            throw PHNetworkingError.invalidUrl
        }

        // TODO: make request
    }
}
