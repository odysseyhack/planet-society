//
//  PermissionFlowContainerViewController.swift
//  Permission Hub
//
//  Created by Corné on 12/04/2019.
//  Copyright © 2019 Planet. All rights reserved.
//

import UIKit

struct PermissionFlowStep {
    let viewController = UIViewController()
}

final class PermissionFlowContainerViewController: UIPageViewController {

    // MARK: - Private properties

    private let steps: [PermissionFlowStep]
    private var currentStepIndex = 0

    private let service = NetworkingService()

    // MARK: - Initialization

    init(steps: [PermissionFlowStep]) {
        self.steps = steps

        super.init(
            transitionStyle: .scroll,
            navigationOrientation: .horizontal,
            options: [:])
    }

    required init?(coder aDecoder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }

    // MARK: - Life cycle

    override func viewDidLoad() {
        super.viewDidLoad()

        view.backgroundColor = Colors.lightGray

        // set initial viewcontroller
        let firstViewControllerOrNil = self.steps.map { $0.viewController }.first
        if let viewController = firstViewControllerOrNil {

            setViewControllers(
                [viewController],
                direction: .forward,
                animated: false)
        }

        pollForNotifications()
    }

    // MARK: - Networking

    private func pollForNotifications() {

        try! service.getNotifications { result in

            switch result {
            case .success(let notification):
                print(notification)
            case .failure(let error):
                print(error)
            }
        }
    }
}
